import type { ChangeEvent, SyntheticEvent } from 'react';
import { useEffect, useState } from 'react';

import { useNavigate } from 'react-router';

import { IconPlus, IconUser } from '@tabler/icons-react';

import {
  Box,
  Button,
  Card,
  Dialog,
  Flex,
  IconButton,
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
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newClientName, setNewClientName] = useState('');
  const [creating, setCreating] = useState(false);
  const { token } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    async function fetchClients() {
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
    }

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
        // Refresh clients list
        const refreshResponse = await fetch('/api/v1/clients', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        if (refreshResponse.ok) {
          const data = await refreshResponse.json();
          setClients(data || []);
        }
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

  return (
    <Stack gap="6">
      <Flex justify="between" align="center">
        <Box>
          <Text size="6" weight="bold">
            Clients
          </Text>
          <Text color="muted">
            Manage your client base and their associated projects.
          </Text>
        </Box>
        <Button variant="filled" onClick={() => setIsModalOpen(true)}>
          <Stack direction="row" gap="2" align="center">
            <IconPlus size={18} />
            <Text>New Client</Text>
          </Stack>
        </Button>
      </Flex>

      <Card p="0" style={{ overflow: 'hidden' }}>
        <Table.Root>
          <Table.Header>
            <Table.Row>
              <Table.ColumnHeader>Client Name</Table.ColumnHeader>
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
                    <Text color="muted">Loading clients...</Text>
                  </Flex>
                </Table.Cell>
              </Table.Row>
            ) : clients.length === 0 ? (
              <Table.Row>
                <Table.Cell colSpan={4}>
                  <Flex p="8" justify="center">
                    <Text color="muted">No clients found.</Text>
                  </Flex>
                </Table.Cell>
              </Table.Row>
            ) : (
              clients.map((client) => (
                <Table.Row key={client.id}>
                  <Table.Cell>
                    <Flex gap="3" align="center">
                      <IconButton
                        variant="text"
                        color="source"
                        style={{
                          backgroundColor:
                            'rgba(var(--pittorica-color-accent-rgb), 0.1)',
                        }}
                      >
                        <IconUser size={16} />
                      </IconButton>
                      <Text weight="medium">{client.name}</Text>
                    </Flex>
                  </Table.Cell>
                  <Table.Cell>
                    <Text
                      size="1"
                      color="muted"
                      style={{ fontFamily: 'var(--pittorica-font-code)' }}
                    >
                      {client.id}
                    </Text>
                  </Table.Cell>
                  <Table.Cell>
                    <Text size="2">
                      {new Date(client.created_at).toLocaleDateString()}
                    </Text>
                  </Table.Cell>
                  <Table.Cell style={{ textAlign: 'right' }}>
                    <Button
                      variant="text"
                      size="sm"
                      onClick={() =>
                        navigate(`/dashboard/clients/${client.id}`)
                      }
                    >
                      View Projects
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
          <Stack gap="4">
            <Text color="muted">
              Enter the name of the new client to add to the vault.
            </Text>
            <TextField.Root>
              <TextField.Input
                placeholder="e.g. Acme Corp"
                autoFocus
                value={newClientName}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewClientName(e.target.value)
                }
              />
            </TextField.Root>
            <Flex justify="end" gap="3">
              <Button variant="text" onClick={() => setIsModalOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" variant="filled" disabled={creating}>
                {creating ? 'Creating...' : 'Create Client'}
              </Button>
            </Flex>
          </Stack>
        </form>
      </Dialog>
    </Stack>
  );
}
