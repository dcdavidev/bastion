import type { ChangeEvent, FormEvent } from 'react';
import { useEffect, useState } from 'react';

import { useNavigate, useParams } from 'react-router';

import {
  IconAlertTriangle,
  IconBriefcase,
  IconChevronLeft,
  IconPlus,
  IconSearch,
  IconShieldLock,
  IconTrash,
} from '@tabler/icons-react';

import {
  Badge,
  Button,
  Card,
  Dialog,
  Divider,
  Flex,
  Stack,
  Table,
  Text,
  TextField,
  toast,
} from '@pittorica/react';

import { useAuth } from '../contexts/auth-context';
import {
  bytesToHex,
  decrypt,
  deriveKey,
  encrypt,
  hexToBytes,
} from '../utils/crypto';

interface Project {
  id: string;
  name: string;
  created_at: string;
  wrapped_data_key: string;
}

interface Client {
  id: string;
  name: string;
}

export default function ClientDetail() {
  const { clientId } = useParams();
  const navigate = useNavigate();
  const [projects, setProjects] = useState<Project[]>([]);
  const [client, setClient] = useState<Client | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(true);

  // Modals
  const [isProjectModalOpen, setIsProjectModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);

  // Form States
  const [newProjectName, setNewProjectName] = useState('');
  const [adminPassword, setAdminPassword] = useState('');
  const [creating, setCreating] = useState(false);
  const [deleting, setDeleting] = useState(false);

  const { token } = useAuth();

  const fetchData = async () => {
    if (!token) return;
    try {
      const [pResp, cResp] = await Promise.all([
        fetch(`/api/v1/projects?client_id=${clientId}`, {
          headers: { Authorization: `Bearer ${token}` },
        }),
        fetch(`/api/v1/clients`, {
          headers: { Authorization: `Bearer ${token}` },
        }),
      ]);

      if (pResp.ok) {
        const pData = await pResp.json();
        setProjects(Array.isArray(pData) ? pData : []);
      }

      if (cResp.ok) {
        const cData = await cResp.json();
        const currentClient = (cData as Client[]).find(
          (c) => c.id === clientId
        );
        if (currentClient) setClient(currentClient);
      }
    } catch (error) {
      console.error('Failed to fetch data', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
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
        setIsProjectModalOpen(false);
        toast({
          title: 'Project secured',
          description: `Environment ${newProjectName} has been created.`,
          color: 'teal',
        });
        fetchData();
      }
    } catch {
      toast({
        title: 'Failed',
        description: 'Check password or network.',
        color: 'red',
      });
    } finally {
      setCreating(false);
    }
  }

  async function handleDeleteClient() {
    if (!token) return;
    setDeleting(true);
    try {
      const response = await fetch(`/api/v1/clients/${clientId}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      });

      if (response.ok) {
        toast({
          title: 'Client deleted',
          description: 'The client and all associated data have been removed.',
          color: 'teal',
        });
        navigate('/vault/clients');
      } else {
        throw new Error('Server error');
      }
    } catch {
      toast({
        title: 'Error',
        description: 'Failed to delete client.',
        color: 'red',
      });
    } finally {
      setDeleting(false);
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
        onClick={() => navigate('/vault/clients')}
        style={{ width: 'fit-content', marginLeft: '-8px' }}
      >
        <Flex gap="2" align="center">
          <IconChevronLeft size={18} />
          <Text>Back to Clients</Text>
        </Flex>
      </Button>

      <Flex gap="6" align="start">
        {/* Main Content */}
        <Stack gap="6" style={{ flex: 1 }}>
          <Stack gap="2">
            <Text size="7" weight="bold" color="source">
              {client?.name || 'Loading...'}
            </Text>
            <Text color="muted" size="2">
              Isolated E2EE environments for this client.
            </Text>
          </Stack>

          <Card p="4">
            <TextField.Root size="md">
              <TextField.Slot>
                <IconSearch size={18} />
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
                  <Table.ColumnHeader>Environment</Table.ColumnHeader>
                  <Table.ColumnHeader>UUID</Table.ColumnHeader>
                  <Table.ColumnHeader style={{ textAlign: 'right' }}>
                    Actions
                  </Table.ColumnHeader>
                </Table.Row>
              </Table.Header>
              <Table.Body>
                {loading ? (
                  <Table.Row>
                    <Table.Cell colSpan={3}>
                      <Flex p="6" justify="center">
                        <Text color="muted">Syncing...</Text>
                      </Flex>
                    </Table.Cell>
                  </Table.Row>
                ) : filteredProjects.length === 0 ? (
                  <Table.Row>
                    <Table.Cell colSpan={3}>
                      <Flex p="8" justify="center">
                        <Text color="muted">No environments found.</Text>
                      </Flex>
                    </Table.Cell>
                  </Table.Row>
                ) : (
                  filteredProjects.map((p) => (
                    <Table.Row key={p.id}>
                      <Table.Cell>
                        <Flex gap="2" align="center">
                          <IconBriefcase
                            size={16}
                            color="var(--pittorica-color-source)"
                          />
                          <Text weight="bold">{p.name}</Text>
                        </Flex>
                      </Table.Cell>
                      <Table.Cell>
                        <Badge variant="standard">
                          <Text
                            size="1"
                            style={{ fontFamily: 'var(--pittorica-font-code)' }}
                          >
                            {p.id.slice(0, 8)}...
                          </Text>
                        </Badge>
                      </Table.Cell>
                      <Table.Cell style={{ textAlign: 'right' }}>
                        <Button
                          variant="outlined"
                          size="sm"
                          onClick={() => navigate(`/vault/projects/${p.id}`)}
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
        </Stack>

        {/* Sidebar */}
        <Stack gap="6" style={{ width: '300px' }}>
          <Card p="5">
            <Stack gap="4">
              <Text weight="bold" size="3">
                Actions
              </Text>
              <Button
                variant="filled"
                style={{ width: '100%' }}
                onClick={() => setIsProjectModalOpen(true)}
              >
                <IconPlus size={18} />
                <Text>New Project</Text>
              </Button>
            </Stack>
          </Card>

          <Card
            p="5"
            style={{
              backgroundColor:
                'rgba(var(--pittorica-color-red-rgb, 244, 67, 54), 0.05)',
              border: '1px solid var(--pittorica-color-red, #f44336)',
            }}
          >
            <Stack gap="4">
              <Text weight="bold" size="3" color="red">
                Danger Zone
              </Text>
              <Text size="1" color="muted">
                Irreversible actions. Deleting a client will permanently remove
                all associated projects and E2EE secrets.
              </Text>
              <Divider />
              <Button
                variant="tonal"
                color="red"
                style={{ width: '100%' }}
                onClick={() => setIsDeleteModalOpen(true)}
              >
                <IconTrash size={18} />
                <Text>Delete Client</Text>
              </Button>
            </Stack>
          </Card>
        </Stack>
      </Flex>

      {/* New Project Modal */}
      <Dialog
        open={isProjectModalOpen}
        onClose={() => setIsProjectModalOpen(false)}
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
                borderRadius: '8px',
              }}
            >
              <IconShieldLock size={24} color="var(--pittorica-color-source)" />
              <Text size="2" color="muted">
                Generates a 256-bit AES key, wrapped with your master key.
              </Text>
            </Flex>
            <TextField.Root size="md" label="Project Name">
              <TextField.Input
                placeholder="e.g. Production"
                value={newProjectName}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewProjectName(e.target.value)
                }
                required
              />
            </TextField.Root>
            <TextField.Root size="md" label="Master Password">
              <TextField.Input
                type="password"
                placeholder="Required to wrap key"
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
                onClick={() => setIsProjectModalOpen(false)}
              >
                Cancel
              </Button>
              <Button type="submit" variant="filled" disabled={creating}>
                {creating ? 'Securing...' : 'Create & Secure'}
              </Button>
            </Flex>
          </Stack>
        </form>
      </Dialog>

      {/* Delete Client Modal */}
      <Dialog
        open={isDeleteModalOpen}
        onClose={() => setIsDeleteModalOpen(false)}
        title="Delete Client?"
      >
        <Stack gap="5">
          <Flex
            gap="3"
            align="center"
            p="3"
            style={{
              backgroundColor: 'rgba(var(--pittorica-color-red-rgb), 0.05)',
              borderRadius: '8px',
            }}
          >
            <IconAlertTriangle size={24} color="var(--pittorica-color-red)" />
            <Text size="2" color="red">
              Warning: This action is permanent and will destroy all secrets for
              this tenant.
            </Text>
          </Flex>
          <Text size="2">
            Are you sure you want to delete <strong>{client?.name}</strong>? All
            projects and encrypted data will be wiped from the vault.
          </Text>
          <Flex justify="end" gap="3">
            <Button variant="text" onClick={() => setIsDeleteModalOpen(false)}>
              Keep Client
            </Button>
            <Button
              variant="filled"
              color="red"
              onClick={handleDeleteClient}
              disabled={deleting}
            >
              {deleting ? 'Deleting...' : 'Yes, Delete Everything'}
            </Button>
          </Flex>
        </Stack>
      </Dialog>
    </Stack>
  );
}
