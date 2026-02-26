import type { ChangeEvent, FormEvent } from 'react';
import { useEffect, useState } from 'react';

import { useNavigate, useParams } from 'react-router';

import {
  IconBriefcase,
  IconChevronLeft,
  IconPlus,
  IconSearch,
  IconShieldLock,
} from '@tabler/icons-react';

import {
  Badge,
  Box,
  Button,
  Card,
  Dialog,
  Flex,
  Stack,
  Table,
  Text,
  TextField,
  toast,
} from '@pittorica/react';

import { useAuth } from '../../contexts/auth-context';
import {
  bytesToHex,
  decrypt,
  deriveKey,
  encrypt,
  hexToBytes,
} from '../../utils/crypto';

interface Project {
  id: string;
  name: string;
  created_at: string;
  wrapped_data_key: string;
}

export default function Projects() {
  const { clientId } = useParams();
  const navigate = useNavigate();
  const [projects, setProjects] = useState<Project[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newProjectName, setNewProjectName] = useState('');
  const [adminPassword, setAdminPassword] = useState('');
  const [creating, setCreating] = useState(false);
  const { token } = useAuth();

  const fetchProjects = async () => {
    try {
      const response = await fetch(`/api/v1/projects?client_id=${clientId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setProjects(data || []);
      }
    } catch (error) {
      console.error('Failed to fetch projects', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProjects();
  }, [clientId, token]);

  async function handleCreateProject(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!newProjectName.trim() || !adminPassword) return;
    setCreating(true);

    try {
      const vcResponse = await fetch('/api/v1/vault/config', {
        headers: { Authorization: `Bearer ${token}` },
      });
      const vc = await vcResponse.json();

      const salt = hexToBytes(vc.master_key_salt);
      const wrappedMK = hexToBytes(vc.wrapped_master_key);
      const adminKEK = await deriveKey(adminPassword, salt);
      const masterKey = await decrypt(adminKEK, wrappedMK);

      const dataKey = globalThis.crypto.getRandomValues(new Uint8Array(32));
      const wrappedDK = await encrypt(masterKey, dataKey);

      const response = await fetch('/api/v1/projects', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          client_id: clientId,
          name: newProjectName,
          wrapped_data_key: bytesToHex(wrappedDK),
        }),
      });

      if (response.ok) {
        setNewProjectName('');
        setAdminPassword('');
        setIsModalOpen(false);
        toast({
          title: 'Project secured',
          description: `Environment ${newProjectName} has been created with a unique E2EE key.`,
          color: 'teal',
        });
        fetchProjects();
      } else {
        throw new Error('Failed to create project on server');
      }
    } catch (error: unknown) {
      console.error('Encryption or request failed', error);
      toast({
        title: 'Project creation failed',
        description:
          error instanceof Error
            ? error.message
            : 'Invalid admin password or crypto error',
        color: 'red',
      });
    } finally {
      setCreating(false);
    }
  }

  const filteredProjects = projects.filter(
    (p) =>
      p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.id.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <Stack gap="6">
      <Button
        variant="text"
        size="sm"
        onClick={() => navigate('/dashboard/clients')}
        style={{ width: 'fit-content', marginLeft: '-8px' }}
      >
        <Flex gap="2" align="center">
          <IconChevronLeft size={18} />
          <Text>Back to Clients</Text>
        </Flex>
      </Button>

      <Flex justify="between" align="end">
        <Box>
          <Text size="7" weight="bold" color="source">
            Projects
          </Text>
          <Text color="muted" size="2">
            Isolated E2EE environments for this client.
          </Text>
        </Box>
        <Button variant="filled" size="md" onClick={() => setIsModalOpen(true)}>
          <Flex gap="2" align="center">
            <IconPlus size={18} />
            <Text>New Project</Text>
          </Flex>
        </Button>
      </Flex>

      <Card p="4">
        <TextField.Root size="md">
          <TextField.Slot>
            <IconSearch size={18} color="var(--pittorica-color-muted)" />
          </TextField.Slot>
          <TextField.Input
            placeholder="Search projects..."
            value={searchQuery}
            onChange={(e: ChangeEvent<HTMLInputElement>) =>
              setSearchQuery(e.target.value)
            }
          />
        </TextField.Root>
      </Card>

      <Card p="0" style={{ overflow: 'hidden' }}>
        <Table.Root>
          <Table.Header>
            <Table.Row>
              <Table.ColumnHeader>Environment Name</Table.ColumnHeader>
              <Table.ColumnHeader>UUID</Table.ColumnHeader>
              <Table.ColumnHeader>Created At</Table.ColumnHeader>
              <Table.ColumnHeader style={{ textAlign: 'right' }}>
                Actions
              </Table.ColumnHeader>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {loading ? (
              <Table.Row>
                <Table.Cell colSpan={4}>
                  <Flex p="8" justify="center">
                    <Text color="muted">Loading projects...</Text>
                  </Flex>
                </Table.Cell>
              </Table.Row>
            ) : filteredProjects.length === 0 ? (
              <Table.Row>
                <Table.Cell colSpan={4}>
                  <Flex p="8" justify="center">
                    <Stack align="center" gap="2">
                      <IconBriefcase
                        size={32}
                        color="var(--pittorica-color-muted)"
                      />
                      <Text color="muted">No environments found.</Text>
                    </Stack>
                  </Flex>
                </Table.Cell>
              </Table.Row>
            ) : (
              filteredProjects.map((p) => (
                <Table.Row key={p.id}>
                  <Table.Cell>
                    <Flex gap="3" align="center">
                      <Box
                        p="2"
                        style={{
                          backgroundColor:
                            'rgba(var(--pittorica-color-source-rgb), 0.1)',
                          borderRadius: 'var(--pittorica-radius-full)',
                        }}
                      >
                        <IconBriefcase
                          size={16}
                          color="var(--pittorica-color-source)"
                        />
                      </Box>
                      <Text weight="bold">{p.name}</Text>
                    </Flex>
                  </Table.Cell>
                  <Table.Cell>
                    <Badge variant="standard">
                      <Text
                        style={{ fontFamily: 'var(--pittorica-font-code)' }}
                      >
                        {p.id}
                      </Text>
                    </Badge>
                  </Table.Cell>
                  <Table.Cell>
                    <Text size="2" color="muted">
                      {new Date(p.created_at).toLocaleDateString(undefined, {
                        dateStyle: 'medium',
                      })}
                    </Text>
                  </Table.Cell>
                  <Table.Cell style={{ textAlign: 'right' }}>
                    <Button
                      variant="outlined"
                      size="sm"
                      onClick={() => navigate(`/dashboard/projects/${p.id}`)}
                    >
                      Manage Secrets
                    </Button>
                  </Table.Cell>
                </Table.Row>
              ))
            )}
          </Table.Body>
        </Table.Root>
      </Card>

      <Dialog
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Secure New Environment"
      >
        <form onSubmit={handleCreateProject}>
          <Stack gap="5">
            <Flex
              gap="3"
              align="center"
              p="3"
              style={{
                backgroundColor:
                  'rgba(var(--pittorica-color-source-rgb), 0.05)',
                borderRadius: 'var(--pittorica-radius-md)',
              }}
            >
              <IconShieldLock size={24} color="var(--pittorica-color-source)" />
              <Text size="2" color="muted">
                This will generate a 256-bit AES key, wrapped with the vault
                master key.
              </Text>
            </Flex>

            <TextField.Root size="md" label="Project Name">
              <TextField.Input
                placeholder="e.g. Production, Staging"
                value={newProjectName}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewProjectName(e.target.value)
                }
                required
              />
            </TextField.Root>
            <TextField.Root size="md" label="Confirm Identity">
              <TextField.Input
                type="password"
                placeholder="Admin Password"
                value={adminPassword}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setAdminPassword(e.target.value)
                }
                required
              />
            </TextField.Root>
            <Flex justify="end" gap="3">
              <Button
                variant="text"
                size="md"
                onClick={() => setIsModalOpen(false)}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                variant="filled"
                size="md"
                disabled={creating || !newProjectName.trim() || !adminPassword}
              >
                {creating ? 'Processing...' : 'Create & Secure'}
              </Button>
            </Flex>
          </Stack>
        </form>
      </Dialog>
    </Stack>
  );
}
