import type { ChangeEvent, SyntheticEvent } from 'react';
import { useEffect, useState } from 'react';

import { useNavigate } from 'react-router';

import {
  IconExternalLink,
  IconPlus,
  IconSearch,
  IconUser,
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

interface Client {
  id: string;
  name: string;
  created_at: string;
}

export default function Clients() {
  const [clients, setClients] = useState<Client[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newClientName, setNewClientName] = useState('');
  const [creating, setCreating] = useState(false);
  const { token } = useAuth();
  const navigate = useNavigate();

  const fetchClients = async () => {
    try {
      const response = await fetch('/api/v1/clients', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (response.ok) {
        const data = await response.json();
        setClients(data || []);
      }
    } catch (error) {
      console.error('Failed to fetch clients', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchClients();
  }, [token]);

  async function handleCreateClient(e: SyntheticEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!newClientName.trim()) return;

    setCreating(true);
    try {
      const response = await fetch('/api/v1/clients', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ name: newClientName }),
      });

      if (response.ok) {
        setNewClientName('');
        setIsModalOpen(false);
        toast({
          title: 'Client registered',
          description: `Successfully added ${newClientName} to the vault.`,
          color: 'teal',
        });
        fetchClients();
      } else {
        throw new Error('Failed to create client on server');
      }
    } catch (error) {
      console.error('Failed to create client', error);
      toast({
        title: 'Registration failed',
        description:
          error instanceof Error ? error.message : 'An unknown error occurred',
        color: 'red',
      });
    } finally {
      setCreating(false);
    }
  }

  const filteredClients = clients.filter(
    (c) =>
      c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      c.id.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <Stack gap="6">
      <Flex justify="between" align="end">
        <Box>
          <Text size="7" weight="bold" color="source">
            Clients
          </Text>
          <Text color="muted" size="2">
            Manage your client base and their associated projects.
          </Text>
        </Box>
        <Button variant="tonal" size="md" onClick={() => setIsModalOpen(true)}>
          <IconPlus size={18} />
          <Text>New Client</Text>
        </Button>
      </Flex>

      <Card p="4">
        <TextField.Root size="md">
          <TextField.Slot>
            <IconSearch size={18} />
          </TextField.Slot>
          <TextField.Input
            placeholder="Search clients by name or ID..."
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
              <Table.ColumnHeader>Client Name</Table.ColumnHeader>
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
                    <Text color="muted">Loading clients...</Text>
                  </Flex>
                </Table.Cell>
              </Table.Row>
            ) : filteredClients.length === 0 ? (
              <Table.Row>
                <Table.Cell colSpan={4}>
                  <Flex p="8" justify="center">
                    <Stack align="center" gap="2">
                      <IconUser
                        size={32}
                        color="var(--pittorica-color-muted)"
                      />
                      <Text color="muted">No clients found.</Text>
                    </Stack>
                  </Flex>
                </Table.Cell>
              </Table.Row>
            ) : (
              filteredClients.map((client) => (
                <Table.Row key={client.id}>
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
                        <IconUser
                          size={16}
                          color="var(--pittorica-color-source)"
                        />
                      </Box>
                      <Text weight="bold">{client.name}</Text>
                    </Flex>
                  </Table.Cell>
                  <Table.Cell>
                    <Badge variant="standard">
                      <Text
                        style={{ fontFamily: 'var(--pittorica-font-code)' }}
                      >
                        {client.id}
                      </Text>
                    </Badge>
                  </Table.Cell>
                  <Table.Cell>
                    <Text size="2" color="muted">
                      {new Date(client.created_at).toLocaleDateString(
                        undefined,
                        { dateStyle: 'medium' }
                      )}
                    </Text>
                  </Table.Cell>
                  <Table.Cell style={{ textAlign: 'right' }}>
                    <Button
                      variant="tonal"
                      size="sm"
                      onClick={() =>
                        navigate(`/dashboard/clients/${client.id}`)
                      }
                    >
                      <Flex gap="1" align="center">
                        <Text>View Projects</Text>
                        <IconExternalLink size={14} />
                      </Flex>
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
        title="Register New Client"
      >
        <form onSubmit={handleCreateClient}>
          <Stack gap="5">
            <Text color="muted" size="2">
              Add a new tenant to the vault. Each client can manage multiple
              isolated projects.
            </Text>
            <TextField.Root size="md" label="Client Name">
              <TextField.Input
                placeholder="e.g. Acme Corporation"
                autoFocus
                value={newClientName}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewClientName(e.target.value)
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
                disabled={creating || !newClientName.trim()}
              >
                {creating ? 'Creating...' : 'Register Client'}
              </Button>
            </Flex>
          </Stack>
        </form>
      </Dialog>
    </Stack>
  );
}
