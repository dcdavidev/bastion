import type { ChangeEvent, FormEvent } from 'react';
import { useState } from 'react';

import {
  IconAlertCircle,
  IconShield,
  IconUserPlus,
  IconUsers,
} from '@tabler/icons-react';

import {
  Box,
  Button,
  Card,
  Chip,
  Dialog,
  Flex,
  Grid,
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
    email: '',
    password: '',
    projectId: '',
  });
  const [creating, setCreating] = useState(false);
  const { token } = useAuth();

  async function handleCreateCollaborator(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (
      !newCollab.username ||
      !newCollab.email ||
      !newCollab.password ||
      !newCollab.projectId
    )
      return;

    setCreating(true);

    // In a real scenario, we would derive the KEK client-side
    // For now, keeping the mock logic as per backend requirements
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
          email: newCollab.email,
          password_hash: 'mock-hash', // Should be Argon2id in production
          salt: 'mock-salt',
          project_id: newCollab.projectId,
          wrapped_data_key: mockWrappedKey,
        }),
      });

      if (response.ok) {
        setIsModalOpen(false);
        const addedUsername = newCollab.username;
        setNewCollab({ username: '', email: '', password: '', projectId: '' });
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
      <Flex justify="between" align="end">
        <Box>
          <Text size="7" weight="bold" color="source">
            Collaborators
          </Text>
          <Text color="muted" size="2">
            Manage restricted access for team members.
          </Text>
        </Box>
        <Button variant="filled" size="md" onClick={() => setIsModalOpen(true)}>
          <Flex gap="2" align="center">
            <IconUserPlus size={18} />
            <Text>New Collaborator</Text>
          </Flex>
        </Button>
      </Flex>

      <Grid columns="2" gap="6">
        <Card p="6">
          <Stack gap="4">
            <Flex color="source" align="center" gap="3">
              <IconShield size={32} />
              <Text size="4" weight="bold">
                Access Control Policy
              </Text>
            </Flex>
            <Text color="muted" size="2">
              Collaborators are restricted users who only possess the keys to
              specific projects. They cannot create clients or view audit logs
              unless explicitly authorized.
            </Text>
            <Box
              p="3"
              style={{
                backgroundColor:
                  'rgba(var(--pittorica-color-source-rgb), 0.05)',
                borderRadius: 'var(--pittorica-radius-md)',
              }}
            >
              <Flex gap="2" align="start">
                <IconAlertCircle
                  size={18}
                  color="var(--pittorica-color-source)"
                  style={{ marginTop: '2px' }}
                />
                <Text size="1" color="muted">
                  When you add a collaborator, a new unique re-wrapped project
                  key is generated specifically for their password.
                </Text>
              </Flex>
            </Box>
          </Stack>
        </Card>

        <Card p="6">
          <Flex
            direction={'column'}
            gap="4"
            align="center"
            justify="center"
            style={{ height: '100%', textAlign: 'center' }}
          >
            <IconUsers size={48} color="var(--pittorica-color-muted)" />
            <Stack>
              <Text weight="bold">No active sessions</Text>
              <Text size="2" color="muted">
                Collaborator management is in read-only mode.
              </Text>
            </Stack>
            <Chip variant="soft">Coming Soon: Active Sessions List</Chip>
          </Flex>
        </Card>
      </Grid>

      <Dialog
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Add Project Collaborator"
      >
        <form onSubmit={handleCreateCollaborator}>
          <Stack gap="5">
            <Text color="muted" size="2">
              Assign a user to a specific project with restricted credentials.
            </Text>
            <Grid columns="2" gap="4">
              <TextField.Root size="md" label="Username">
                <TextField.Input
                  placeholder="John Doe"
                  value={newCollab.username}
                  onChange={(e: ChangeEvent<HTMLInputElement>) =>
                    setNewCollab({ ...newCollab, username: e.target.value })
                  }
                  required
                />
              </TextField.Root>
              <TextField.Root size="md" label="Email Address">
                <TextField.Input
                  type="email"
                  placeholder="john@example.com"
                  value={newCollab.email}
                  onChange={(e: ChangeEvent<HTMLInputElement>) =>
                    setNewCollab({ ...newCollab, email: e.target.value })
                  }
                  required
                />
              </TextField.Root>
            </Grid>
            <TextField.Root size="md" label="Temporary Password">
              <TextField.Input
                type="password"
                placeholder="Assign a secure password"
                value={newCollab.password}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewCollab({ ...newCollab, password: e.target.value })
                }
                required
              />
            </TextField.Root>
            <TextField.Root size="md" label="Target Project UUID">
              <TextField.Input
                placeholder="00000000-0000-0000-0000-000000000000"
                value={newCollab.projectId}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewCollab({ ...newCollab, projectId: e.target.value })
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
                disabled={creating}
              >
                {creating ? 'Provisioning...' : 'Grant Access'}
              </Button>
            </Flex>
          </Stack>
        </form>
      </Dialog>
    </Stack>
  );
}
