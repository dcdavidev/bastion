import type { ChangeEvent, FormEvent } from 'react';
import { useEffect, useState } from 'react';

import { useNavigate, useParams } from 'react-router';

import {
  IconChevronLeft,
  IconEye,
  IconEyeOff,
  IconLock,
  IconLockOpen,
  IconPlus,
  IconSearch,
} from '@tabler/icons-react';

import {
  Badge,
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
import {
  bytesToHex,
  decrypt,
  deriveKey,
  encrypt,
  hexToBytes,
} from '../../utils/crypto';

interface Secret {
  id: string;
  key: string;
  value: string;
  version: number;
}

export default function Secrets() {
  const { projectId } = useParams();
  const navigate = useNavigate();
  const [secrets, setSecrets] = useState<Secret[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [decryptedValues, setDecryptedValues] = useState<
    Record<string, string>
  >({});
  const [visibleSecrets, setVisibleSecrets] = useState<Record<string, boolean>>(
    {}
  );

  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newSecret, setNewSecret] = useState({ key: '', value: '' });
  const [adminPassword, setAdminPassword] = useState('');
  const [unlocking, setUnlocking] = useState(false);
  const [unlocked, setUnlocked] = useState(false);
  const [projectDataKey, setProjectDataKey] = useState<Uint8Array | null>(null);
  const [creating, setCreating] = useState(false);
  const { token } = useAuth();

  const fetchSecrets = async () => {
    try {
      const response = await fetch(`/api/v1/secrets?project_id=${projectId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setSecrets(data || []);

        // If already unlocked, re-decrypt new data
        if (projectDataKey) {
          const newDecrypted: Record<string, string> = {};
          for (const s of data as Secret[]) {
            try {
              const ciphertext = hexToBytes(s.value);
              const plaintext = await decrypt(projectDataKey, ciphertext);
              newDecrypted[s.id] = new TextDecoder().decode(plaintext);
            } catch {
              console.error('Failed to decrypt', s.id);
            }
          }
          setDecryptedValues(newDecrypted);
        }
      }
    } catch (error) {
      console.error('Failed to fetch secrets', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSecrets();
  }, [projectId, token]);

  async function handleUnlock(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!adminPassword) return;
    setUnlocking(true);
    try {
      const [vcResponse, pResponse] = await Promise.all([
        fetch('/api/v1/vault/config', {
          headers: { Authorization: `Bearer ${token}` },
        }),
        fetch(`/api/v1/projects/${projectId}`, {
          headers: { Authorization: `Bearer ${token}` },
        }),
      ]);

      const vc = await vcResponse.json();
      const project = await pResponse.json();

      const salt = hexToBytes(vc.master_key_salt);
      const wrappedMK = hexToBytes(vc.wrapped_master_key);
      const adminKEK = await deriveKey(adminPassword, salt);
      const masterKey = await decrypt(adminKEK, wrappedMK);

      const wrappedDK = hexToBytes(project.wrapped_data_key);
      const dataKey = await decrypt(masterKey, wrappedDK);

      setProjectDataKey(dataKey);
      setUnlocked(true);

      // Decrypt all
      const newDecrypted: Record<string, string> = {};
      for (const s of secrets) {
        try {
          const ciphertext = hexToBytes(s.value);
          const plaintext = await decrypt(dataKey, ciphertext);
          newDecrypted[s.id] = new TextDecoder().decode(plaintext);
        } catch {
          console.error('Failed to decrypt', s.id);
        }
      }
      setDecryptedValues(newDecrypted);

      toast({
        title: 'Vault unlocked',
        description: `${secrets.length} secrets decrypted successfully.`,
        color: 'teal',
      });
    } catch {
      toast({
        title: 'Unlock failed',
        description: 'Check your master password.',
        color: 'red',
      });
    } finally {
      setUnlocking(false);
    }
  }

  async function handleAddSecret(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!projectDataKey || !newSecret.key || !newSecret.value) return;

    setCreating(true);
    try {
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
          key: newSecret.key,
          value: bytesToHex(ciphertext),
        }),
      });

      if (response.ok) {
        setNewSecret({ key: '', value: '' });
        setIsModalOpen(false);
        toast({
          title: 'Secret stored',
          description: `${newSecret.key} is now secured.`,
          color: 'teal',
        });
        await fetchSecrets();
      } else {
        throw new Error('Server error');
      }
    } catch {
      toast({
        title: 'Storage failed',
        description: 'Could not encrypt or save the secret.',
        color: 'red',
      });
    } finally {
      setCreating(false);
    }
  }

  const filteredSecrets = secrets.filter((s) =>
    s.key.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const toggleVisibility = (id: string) => {
    setVisibleSecrets((prev) => ({ ...prev, [id]: !prev[id] }));
  };

  return (
    <Stack gap="6">
      <Button
        variant="text"
        size="sm"
        onClick={() => navigate(-1)}
        style={{ width: 'fit-content', marginLeft: '-8px' }}
      >
        <Flex gap="2" align="center">
          <IconChevronLeft size={18} />
          <Text>Back</Text>
        </Flex>
      </Button>

      {unlocked ? (
        <>
          <Flex justify="between" align="end">
            <Box>
              <Flex align="center" gap="2">
                <IconLockOpen size={24} color="var(--pittorica-color-teal)" />
                <Text size="7" weight="bold" color="source">
                  Secrets
                </Text>
              </Flex>
              <Text color="muted" size="2">
                E2EE secrets managed with AES-256 GCM.
              </Text>
            </Box>
            <Button
              variant="filled"
              size="md"
              onClick={() => setIsModalOpen(true)}
            >
              <IconPlus size={18} />
              <Text>Add Secret</Text>
            </Button>
          </Flex>

          <Card p="4">
            <TextField.Root size="md">
              <TextField.Slot>
                <IconSearch size={18} color="var(--pittorica-color-muted)" />
              </TextField.Slot>
              <TextField.Input
                placeholder="Search keys..."
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
                  <Table.ColumnHeader>Environment Variable</Table.ColumnHeader>
                  <Table.ColumnHeader>Value</Table.ColumnHeader>
                  <Table.ColumnHeader>Version</Table.ColumnHeader>
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
                        <Text color="muted">Syncing...</Text>
                      </Flex>
                    </Table.Cell>
                  </Table.Row>
                ) : filteredSecrets.length === 0 ? (
                  <Table.Row>
                    <Table.Cell colSpan={4}>
                      <Flex p="8" justify="center">
                        <Text color="muted">
                          No secrets in this environment.
                        </Text>
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
                              backgroundColor:
                                'var(--pittorica-color-surface-container)',
                              borderRadius: 'var(--pittorica-radius-sm)',
                              flex: 1,
                              minWidth: '200px',
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
                                ? decryptedValues[s.id] || 'N/A'
                                : '••••••••••••••••'}
                            </Text>
                          </Box>
                          <IconButton
                            variant="text"
                            size="2"
                            onClick={() => toggleVisibility(s.id)}
                          >
                            {visibleSecrets[s.id] ? (
                              <IconEyeOff size={16} />
                            ) : (
                              <IconEye size={16} />
                            )}
                          </IconButton>
                        </Flex>
                      </Table.Cell>
                      <Table.Cell>
                        <Badge variant="standard">v{s.version}</Badge>
                      </Table.Cell>
                      <Table.Cell style={{ textAlign: 'right' }}>
                        <Button variant="text" size="sm" color="red">
                          Delete
                        </Button>
                      </Table.Cell>
                    </Table.Row>
                  ))
                )}
              </Table.Body>
            </Table.Root>
          </Card>
        </>
      ) : (
        <Card p="8" style={{ maxWidth: '480px', margin: '40px auto' }}>
          <form onSubmit={handleUnlock}>
            <Stack gap="6" align="center">
              <Box
                p="4"
                style={{
                  backgroundColor:
                    'rgba(var(--pittorica-color-source-rgb), 0.1)',
                  borderRadius: 'var(--pittorica-radius-full)',
                }}
              >
                <IconLock size={48} color="var(--pittorica-color-source)" />
              </Box>
              <Stack gap="1" align="center" style={{ textAlign: 'center' }}>
                <Text size="6" weight="bold">
                  Protected Environment
                </Text>
                <Text color="muted" size="2">
                  This project is encrypted. Enter your master password to
                  derive the key and access secrets.
                </Text>
              </Stack>
              <TextField.Root size="md" style={{ width: '100%' }}>
                <TextField.Input
                  type="password"
                  placeholder="Master Password"
                  autoFocus
                  value={adminPassword}
                  onChange={(e: ChangeEvent<HTMLInputElement>) =>
                    setAdminPassword(e.target.value)
                  }
                  required
                />
              </TextField.Root>
              <Button
                type="submit"
                variant="filled"
                size="md"
                disabled={unlocking || !adminPassword}
                style={{ width: '100%' }}
              >
                {unlocking ? 'Unwrapping Keys...' : 'Unlock Project'}
              </Button>
            </Stack>
          </form>
        </Card>
      )}

      <Dialog
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
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
            <TextField.Root size="md" label="Plaintext Value">
              <TextField.Input
                placeholder="Enter secret value..."
                value={newSecret.value}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewSecret({ ...newSecret, value: e.target.value })
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
                disabled={creating || !newSecret.key || !newSecret.value}
              >
                {creating ? 'Encrypting...' : 'Encrypt & Save'}
              </Button>
            </Flex>
          </Stack>
        </form>
      </Dialog>
    </Stack>
  );
}
