import type { ChangeEvent, FormEvent } from 'react';
import { useState } from 'react';

import { IconShield, IconUserPlus } from '@tabler/icons-react';

import {
  Box,
  Button,
  Card,
  Dialog,
  Flex,
  Stack,
  Text,
  TextField,
  toast,
} from '@pittorica/react';

import { useAuth } from '../../contexts/auth-context';

export default function Collaborators() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newCollab, setNewCollab] = useState({
    username: '',
    password: '',
    projectId: '',
  });
  const [creating, setCreating] = useState(false);
  const { token } = useAuth();

  async function handleCreateCollaborator(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setCreating(true);

    const mockWrappedKey = 'collab-wrapped-key-' + Math.random().toString(16);

    try {
      const response = await fetch('/api/v1/collaborators', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          username: newCollab.username,
          password_hash: 'mock-hash', // Placeholder
          salt: 'mock-salt', // Placeholder
          project_id: newCollab.projectId,
          wrapped_data_key: mockWrappedKey,
        }),
      });

      if (response.ok) {
        setIsModalOpen(false);
        const addedUsername = newCollab.username;
        setNewCollab({ username: '', password: '', projectId: '' });
        toast({
          title: 'Access granted',
          description: `Collaborator ${addedUsername} has been assigned to the project.`,
          color: 'teal',
        });
      } else {
        throw new Error('Failed to create collaborator on server');
      }
    } catch (error) {
      console.error('Failed to create collaborator', error);
      toast({
        title: 'Assignment failed',
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
            Collaborators
          </Text>
          <Text color="muted">Manage restricted access for team members.</Text>
        </Box>
        <Button variant="filled" onClick={() => setIsModalOpen(true)}>
          <Stack direction="row" gap="2" align="center">
            <IconUserPlus size={18} />
            <Text>New Collaborator</Text>
          </Stack>
        </Button>
      </Flex>

      <Card p="8">
        <Stack gap="4" align="center" style={{ textAlign: 'center' }}>
          <Box color="cyan">
            <IconShield size={48} />
          </Box>
          <Box>
            <Text size="4" weight="bold">
              Access Control
            </Text>
            <Text color="muted">
              Collaborators can only see and use secrets for projects they are
              assigned to.
            </Text>
          </Box>
        </Stack>
      </Card>

      <Dialog
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Add Collaborator"
      >
        <form onSubmit={handleCreateCollaborator}>
          <Stack gap="4">
            <TextField.Root>
              <TextField.Input
                placeholder="Username"
                value={newCollab.username}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewCollab({ ...newCollab, username: e.target.value })
                }
                required
              />
            </TextField.Root>
            <TextField.Root>
              <TextField.Input
                type="password"
                placeholder="Assign a Password"
                value={newCollab.password}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewCollab({ ...newCollab, password: e.target.value })
                }
                required
              />
            </TextField.Root>
            <TextField.Root>
              <TextField.Input
                placeholder="Project ID (UUID)"
                value={newCollab.projectId}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewCollab({ ...newCollab, projectId: e.target.value })
                }
                required
              />
            </TextField.Root>
            <Flex justify="end" gap="3">
              <Button variant="text" onClick={() => setIsModalOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" variant="filled" disabled={creating}>
                {creating ? 'Granting...' : 'Grant Access'}
              </Button>
            </Flex>
          </Stack>
        </form>
      </Dialog>
    </Stack>
  );
}
