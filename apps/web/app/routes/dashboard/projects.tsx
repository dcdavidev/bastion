import type { ChangeEvent, FormEvent } from 'react';
import { useEffect, useState } from 'react';

import { useNavigate, useParams } from 'react-router';

import { IconBriefcase, IconChevronLeft, IconPlus } from '@tabler/icons-react';

import {
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
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newProjectName, setNewProjectName] = useState('');
  const [adminPassword, setAdminPassword] = useState('');
  const [creating, setCreating] = useState(false);
  const { token } = useAuth();

  useEffect(() => {
    async function fetchProjects() {
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
    }

    fetchProjects();
  }, [clientId, token]);

  async function handleCreateProject(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
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
        // Refresh
        const refreshResponse = await fetch(
          `/api/v1/projects?client_id=${clientId}`,
          {
            headers: { Authorization: `Bearer ${token}` },
          }
        );
        if (refreshResponse.ok) {
          const data = await refreshResponse.json();
          setProjects(data || []);
        }
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

  return (
    <Stack gap="6">
      <Button
        variant="text"
        onClick={() => navigate('/dashboard/clients')}
        p="0"
      >
        <Flex gap="2" align="center">
          <IconChevronLeft size={18} />
          <Text>Back to Clients</Text>
        </Flex>
      </Button>

      <Flex justify="between" align="center">
        <Box>
          <Text size="6" weight="bold">
            Projects
          </Text>
          <Text color="muted">
            Manage isolated environments and secrets for this client.
          </Text>
        </Box>
        <Button variant="filled" onClick={() => setIsModalOpen(true)}>
          <Stack direction="row" gap="2" align="center">
            <IconPlus size={18} />
            <Text>New Project</Text>
          </Stack>
        </Button>
      </Flex>

      <Card p="0" style={{ overflow: 'hidden' }}>
        <Table.Root>
          <Table.Header>
            <Table.Row>
              <Table.ColumnHeader>Project Name</Table.ColumnHeader>
              <Table.ColumnHeader>ID</Table.ColumnHeader>
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
            ) : projects.length === 0 ? (
              <Table.Row>
                <Table.Cell colSpan={4}>
                  <Flex p="8" justify="center">
                    <Text color="muted">No projects found.</Text>
                  </Flex>
                </Table.Cell>
              </Table.Row>
            ) : (
              projects.map((p) => (
                <Table.Row key={p.id}>
                  <Table.Cell>
                    <Flex gap="3" align="center">
                      <Box color="cyan">
                        <IconBriefcase size={16} />
                      </Box>
                      <Text weight="medium">{p.name}</Text>
                    </Flex>
                  </Table.Cell>
                  <Table.Cell>
                    <Text
                      size="1"
                      color="muted"
                      style={{ fontFamily: 'var(--pittorica-font-code)' }}
                    >
                      {p.id}
                    </Text>
                  </Table.Cell>
                  <Table.Cell>
                    <Text size="2">
                      {new Date(p.created_at).toLocaleDateString()}
                    </Text>
                  </Table.Cell>
                  <Table.Cell style={{ textAlign: 'right' }}>
                    <Button
                      variant="text"
                      size="sm"
                      onClick={() => navigate(`/dashboard/projects/${p.id}`)}
                    >
                      View Secrets
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
        title="Create New Project"
      >
        <form onSubmit={handleCreateProject}>
          <Stack gap="4">
            <Text color="muted">
              This will generate a new unique E2EE data key for the project.
            </Text>
            <TextField.Root>
              <TextField.Input
                placeholder="Project Name"
                value={newProjectName}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewProjectName(e.target.value)
                }
                required
              />
            </TextField.Root>
            <TextField.Root>
              <TextField.Input
                type="password"
                placeholder="Admin Password (to wrap key)"
                value={adminPassword}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setAdminPassword(e.target.value)
                }
                required
              />
            </TextField.Root>
            <Flex justify="end" gap="3">
              <Button variant="text" onClick={() => setIsModalOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" variant="filled" disabled={creating}>
                {creating ? 'Creating...' : 'Create & Secure'}
              </Button>
            </Flex>
          </Stack>
        </form>
      </Dialog>
    </Stack>
  );
}
