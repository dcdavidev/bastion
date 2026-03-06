import type { ChangeEvent, FormEvent } from 'react';
import { useCallback, useEffect, useState } from 'react';

import { useNavigate, useParams } from 'react-router';

import {
  IconAlertTriangle,
  IconChevronLeft,
  IconDownload,
  IconEye,
  IconEyeOff,
  IconLock,
  IconLockOpen,
  IconPlus,
  IconSearch,
  IconTrash,
} from '@tabler/icons-react';

import {
  Box,
  Button,
  Card,
  Dialog,
  Divider,
  Flex,
  IconButton,
  Stack,
  Table,
  Text,
  TextArea,
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

interface Secret {
  id: string;
  key: string;
  value: string;
  version: number;
}

interface Project {
  id: string;
  name: string;
  client_id: string;
  wrapped_data_key: string;
}

interface Client {
  id: string;
  name: string;
}

export default function SecretsManagement() {
  const { projectId } = useParams();
  const navigate = useNavigate();
  const { token } = useAuth();

  // State
  const [project, setProject] = useState<Project | null>(null);
  const [client, setClient] = useState<Client | null>(null);
  const [secrets, setSecrets] = useState<Secret[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(true);

  // E2EE State
  const [unlocked, setUnlocked] = useState(false);
  const [projectDataKey, setProjectDataKey] = useState<Uint8Array | null>(null);
  const [decryptedValues, setDecryptedValues] = useState<
    Record<string, string>
  >({});
  const [visibleSecrets, setVisibleSecrets] = useState<Record<string, boolean>>(
    {}
  );

  // Modals
  const [isSecretModalOpen, setIsSecretModalOpen] = useState(false);
  const [isDeleteProjectModalOpen, setIsDeleteProjectModalOpen] =
    useState(false);
  const [isUnlockModalOpen, setIsUnlockModalOpen] = useState(!unlocked);

  // Form States
  const [newSecret, setNewSecret] = useState({ key: '', value: '' });
  const [masterPassword, setMasterPassword] = useState('');
  const [processing, setProcessing] = useState(false);

  const fetchSecrets = useCallback(async () => {
    if (!token || !projectId) return;
    try {
      const response = await fetch(`/api/v1/secrets?project_id=${projectId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setSecrets(Array.isArray(data) ? data : []);
      }
    } catch (error) {
      console.error('Failed to fetch secrets', error);
    } finally {
      setLoading(false);
    }
  }, [token, projectId]);

  const fetchProjectAndClient = useCallback(async () => {
    if (!token || !projectId) return;
    try {
      const pResp = await fetch(`/api/v1/projects/${projectId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (pResp.ok) {
        const pData = (await pResp.json()) as Project;
        setProject(pData);

        // Fetch client to get the name for exporting
        const cResp = await fetch('/api/v1/clients', {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (cResp.ok) {
          const clients = (await cResp.json()) as Client[];
          const currentClient = clients.find((c) => c.id === pData.client_id);
          if (currentClient) setClient(currentClient);
        }
      }
    } catch (error) {
      console.error('Failed to fetch project/client', error);
    }
  }, [token, projectId]);

  useEffect(() => {
    void fetchProjectAndClient();
    void fetchSecrets();
  }, [fetchProjectAndClient, fetchSecrets]);

  // Handle Decryption when secrets or key change
  useEffect(() => {
    if (unlocked && projectDataKey && secrets.length > 0) {
      const decryptAll = async () => {
        const newDecrypted: Record<string, string> = {};
        for (const s of secrets) {
          try {
            const ciphertext = hexToBytes(s.value);
            const plaintext = await decrypt(projectDataKey, ciphertext);
            newDecrypted[s.id] = new TextDecoder().decode(plaintext);
          } catch {
            newDecrypted[s.id] = '[Decryption Error]';
          }
        }
        setDecryptedValues(newDecrypted);
      };
      void decryptAll();
    }
  }, [unlocked, projectDataKey, secrets]);

  async function handleUnlock(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!masterPassword || !project) return;
    setProcessing(true);

    try {
      const vcResponse = await fetch('/api/v1/vault/config', {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!vcResponse.ok)
        throw new Error('Failed to fetch vault configuration');
      const vc = await vcResponse.json();

      const salt = hexToBytes(vc.master_key_salt);
      const wrappedMK = hexToBytes(vc.wrapped_master_key);
      const adminKEK = await deriveKey(masterPassword, salt);
      const masterKey = await decrypt(adminKEK, wrappedMK);

      const wrappedDK = hexToBytes(project.wrapped_data_key);
      const dataKey = await decrypt(masterKey, wrappedDK);

      setProjectDataKey(dataKey);
      setUnlocked(true);
      setIsUnlockModalOpen(false);
      setMasterPassword('');
      toast({
        title: 'Project unlocked',
        description: 'Keys unwrapped successfully.',
        color: 'teal',
      });
    } catch (error) {
      console.error('Unlock error:', error);
      toast({
        title: 'Unlock failed',
        description: 'Invalid master password.',
        color: 'red',
      });
    } finally {
      setProcessing(false);
    }
  }

  async function handleAddSecret(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!projectDataKey || !newSecret.key || !newSecret.value) return;
    setProcessing(true);

    try {
      const finalKey = newSecret.key.toUpperCase();
      const plaintext = new TextEncoder().encode(newSecret.value);
      const ciphertext = await encrypt(projectDataKey, plaintext);

      const response = await fetch('/api/v1/secrets', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          project_id: projectId,
          key: finalKey,
          value: bytesToHex(ciphertext),
        }),
      });

      if (response.ok) {
        setNewSecret({ key: '', value: '' });
        setIsSecretModalOpen(false);
        toast({
          title: 'Secret secured',
          description: `${finalKey} has been stored.`,
          color: 'teal',
        });
        void fetchSecrets();
      } else {
        const errorText = await response.text();
        throw new Error(errorText || 'Server error');
      }
    } catch (error: unknown) {
      console.error('Failed to save secret:', error);
      toast({
        title: 'Storage failed',
        description:
          error instanceof Error
            ? error.message
            : 'Could not encrypt or save the secret.',
        color: 'red',
      });
    } finally {
      setProcessing(false);
    }
  }

  async function handleDeleteSecret(secretId: string) {
    if (!token) return;
    try {
      const response = await fetch(`/api/v1/secrets/${secretId}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        toast({ title: 'Secret removed', color: 'teal' });
        void fetchSecrets();
      }
    } catch {
      toast({ title: 'Error', description: 'Failed to delete.', color: 'red' });
    }
  }

  async function handleDeleteProject() {
    if (!token || !project) return;
    setProcessing(true);
    try {
      const response = await fetch(`/api/v1/projects/${projectId}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        toast({ title: 'Project destroyed', color: 'teal' });
        navigate(`/vault/clients/${project.client_id}`);
      }
    } catch {
      toast({ title: 'Error', color: 'red' });
    } finally {
      setProcessing(false);
    }
  }

  const handleExportSecrets = () => {
    if (!unlocked || secrets.length === 0) return;

    const content = secrets
      .map((s) => `${s.key}=${decryptedValues[s.id] || ''}`)
      .join('\n');

    const clientName =
      client?.name.toLowerCase().replaceAll(/\s+/g, '-') || 'unknown';
    const projectName =
      project?.name.toLowerCase().replaceAll(/\s+/g, '-') || 'default';
    const filename = `${clientName}-${projectName}.txt`;

    const blob = new Blob([content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    document.body.append(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(url);

    toast({
      title: 'Secrets exported',
      description: `Saved as ${filename}`,
      color: 'teal',
    });
  };

  const filteredSecrets = secrets.filter((s) =>
    s.key.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <Stack gap="6">
      <Button
        variant="text"
        size="sm"
        onClick={() =>
          project && navigate(`/vault/clients/${project.client_id}`)
        }
        style={{ width: 'fit-content', marginLeft: '-8px' }}
      >
        <Flex gap="2" align="center">
          <IconChevronLeft size={18} />
          <Text>Back to Client</Text>
        </Flex>
      </Button>

      <Flex gap="6" align="start">
        {/* Main Content */}
        <Stack gap="6" style={{ flex: 1 }}>
          <Stack gap="2">
            <Flex align="center" gap="3">
              {unlocked ? (
                <IconLockOpen size={28} color="var(--pittorica-color-teal)" />
              ) : (
                <IconLock size={28} color="var(--pittorica-color-muted)" />
              )}
              <Text size="7" weight="bold" color="source">
                {project?.name || 'Loading...'}
              </Text>
            </Flex>
            <Text color="muted" size="2">
              End-to-end encrypted environment variables.
            </Text>
          </Stack>

          <Card p="4">
            <TextField.Root size="md">
              <TextField.Slot>
                <IconSearch size={18} />
              </TextField.Slot>
              <TextField.Input
                placeholder="Search variables..."
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
                  <Table.ColumnHeader>Variable</Table.ColumnHeader>
                  <Table.ColumnHeader>Value</Table.ColumnHeader>
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
                ) : unlocked ? (
                  filteredSecrets.length === 0 ? (
                    <Table.Row>
                      <Table.Cell colSpan={3}>
                        <Flex p="8" justify="center">
                          <Text color="muted">No secrets found.</Text>
                        </Flex>
                      </Table.Cell>
                    </Table.Row>
                  ) : (
                    filteredSecrets.map((s) => (
                      <Table.Row key={s.id}>
                        <Table.Cell>
                          <Text
                            weight="bold"
                            size="2"
                            style={{ fontFamily: 'var(--pittorica-font-code)' }}
                          >
                            {s.key}
                          </Text>
                        </Table.Cell>
                        <Table.Cell>
                          <Flex gap="2" align="center">
                            <Box
                              p="2"
                              style={{
                                backgroundColor: 'var(--pittorica-surface-3)',
                                borderRadius: '6px',
                                flex: 1,
                                minWidth: '240px',
                              }}
                            >
                              <Text
                                size="2"
                                style={{
                                  fontFamily: 'var(--pittorica-font-code)',
                                  wordBreak: 'break-all',
                                }}
                              >
                                {visibleSecrets[s.id]
                                  ? decryptedValues[s.id] || '•••'
                                  : '••••••••••••••••'}
                              </Text>
                            </Box>
                            <IconButton
                              variant="text"
                              size="2"
                              onClick={() =>
                                setVisibleSecrets((prev) => ({
                                  ...prev,
                                  [s.id]: !prev[s.id],
                                }))
                              }
                            >
                              {visibleSecrets[s.id] ? (
                                <IconEyeOff size={16} />
                              ) : (
                                <IconEye size={16} />
                              )}
                            </IconButton>
                          </Flex>
                        </Table.Cell>
                        <Table.Cell style={{ textAlign: 'right' }}>
                          <IconButton
                            color="red"
                            variant="text"
                            size="2"
                            onClick={() => void handleDeleteSecret(s.id)}
                          >
                            <IconTrash size={16} />
                          </IconButton>
                        </Table.Cell>
                      </Table.Row>
                    ))
                  )
                ) : (
                  <Table.Row>
                    <Table.Cell colSpan={3}>
                      <Flex p="8" direction="column" align="center" gap="4">
                        <IconLock
                          size={48}
                          color="var(--pittorica-color-muted)"
                        />
                        <Stack align="center" gap="1">
                          <Text weight="bold">Vault Locked</Text>
                          <Text size="2" color="muted">
                            Unlock the project to view and manage secrets.
                          </Text>
                        </Stack>
                        <Button
                          variant="filled"
                          onClick={() => setIsUnlockModalOpen(true)}
                        >
                          Unlock Now
                        </Button>
                      </Flex>
                    </Table.Cell>
                  </Table.Row>
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
                Management
              </Text>
              <Button
                variant="filled"
                style={{ width: '100%' }}
                disabled={!unlocked}
                onClick={() => setIsSecretModalOpen(true)}
              >
                <IconPlus size={18} />
                <Text>Add Secret</Text>
              </Button>
              {!unlocked && (
                <Button
                  variant="tonal"
                  style={{ width: '100%' }}
                  onClick={() => setIsUnlockModalOpen(true)}
                >
                  <IconLockOpen size={18} />
                  <Text>Unlock Keys</Text>
                </Button>
              )}
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
                Sensitive actions. Exporting will download secrets in plaintext.
                Destroying is permanent.
              </Text>
              <Divider />
              <Button
                variant="tonal"
                color="red"
                style={{ width: '100%' }}
                disabled={!unlocked || secrets.length === 0}
                onClick={handleExportSecrets}
              >
                <IconDownload size={18} />
                <Text>Export Secrets (.txt)</Text>
              </Button>
              <Button
                variant="tonal"
                color="red"
                style={{ width: '100%' }}
                onClick={() => setIsDeleteProjectModalOpen(true)}
              >
                <IconTrash size={18} />
                <Text>Destroy Project</Text>
              </Button>
            </Stack>
          </Card>
        </Stack>
      </Flex>

      {/* Unlock Modal */}
      <Dialog
        open={isUnlockModalOpen}
        onClose={() => setIsUnlockModalOpen(false)}
        title="Unlock Environment"
      >
        <form onSubmit={handleUnlock}>
          <Stack gap="5">
            <Text color="muted" size="2">
              Enter your master password to unwrap the project encryption key.
            </Text>
            <TextField.Root size="md" label="Master Password">
              <TextField.Input
                type="password"
                placeholder="Required for E2EE"
                value={masterPassword}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setMasterPassword(e.target.value)
                }
                required
                autoFocus
              />
            </TextField.Root>
            <Flex justify="end" gap="3">
              <Button
                variant="text"
                onClick={() => setIsUnlockModalOpen(false)}
              >
                Cancel
              </Button>
              <Button type="submit" variant="filled" disabled={processing}>
                {processing ? 'Unwrapping...' : 'Unlock Project'}
              </Button>
            </Flex>
          </Stack>
        </form>
      </Dialog>

      {/* Add Secret Modal */}
      <Dialog
        open={isSecretModalOpen}
        onClose={() => setIsSecretModalOpen(false)}
        title="Add Secure Secret"
      >
        <form onSubmit={handleAddSecret}>
          <Stack gap="5">
            <TextField.Root size="md" label="Variable Name">
              <TextField.Input
                placeholder="e.g. DATABASE_URL"
                autoFocus
                value={newSecret.key}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewSecret({
                    ...newSecret,
                    key: e.target.value.toUpperCase(),
                  })
                }
                required
              />
            </TextField.Root>
            <TextArea.Root label="Plaintext Value">
              <TextArea.Content
                placeholder="Enter secret value..."
                value={newSecret.value}
                onChange={(e: ChangeEvent<HTMLTextAreaElement>) =>
                  setNewSecret({ ...newSecret, value: e.target.value })
                }
                required
                rows={4}
                autoResize
              />
            </TextArea.Root>
            <Flex justify="end" gap="3">
              <Button
                variant="text"
                onClick={() => setIsSecretModalOpen(false)}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                variant="filled"
                disabled={processing || !newSecret.key || !newSecret.value}
              >
                {processing ? 'Encrypting...' : 'Encrypt & Save'}
              </Button>
            </Flex>
          </Stack>
        </form>
      </Dialog>

      {/* Delete Project Modal */}
      <Dialog
        open={isDeleteProjectModalOpen}
        onClose={() => setIsDeleteProjectModalOpen(false)}
        title="Destroy Project?"
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
              Warning: This action is permanent.
            </Text>
          </Flex>
          <Text size="2">
            Are you sure you want to destroy <strong>{project?.name}</strong>?
            All secrets will be lost.
          </Text>
          <Flex justify="end" gap="3">
            <Button
              variant="text"
              onClick={() => setIsDeleteProjectModalOpen(false)}
            >
              Cancel
            </Button>
            <Button
              variant="filled"
              color="red"
              onClick={() => void handleDeleteProject()}
              disabled={processing}
            >
              {processing ? 'Destroying...' : 'Yes, Destroy'}
            </Button>
          </Flex>
        </Stack>
      </Dialog>
    </Stack>
  );
}
