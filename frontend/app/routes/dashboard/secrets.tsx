import type { ChangeEvent, FormEvent } from 'react';
import { useEffect, useState } from 'react';

import { useNavigate, useParams } from 'react-router';

import { IconChevronLeft, IconKey, IconPlus } from '@tabler/icons-react';

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
  const [decryptedValues, setDecryptedValues] = useState<
    Record<string, string>
  >({});
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newSecret, setNewSecret] = useState({ key: '', value: '' });
  const [adminPassword, setAdminPassword] = useState('');
  const [unlocking, setUnlocking] = useState(false);
  const [unlocked, setUnlocked] = useState(false);
  const [projectDataKey, setProjectDataKey] = useState<Uint8Array | null>(null);
  const [creating, setCreating] = useState(false);
  const { token } = useAuth();

  useEffect(() => {
    async function fetchSecrets() {
      try {
        const response = await fetch(
          `/api/v1/secrets?project_id=${projectId}`,
          {
            headers: { Authorization: `Bearer ${token}` },
          }
        );
        if (response.ok) {
          const data = await response.json();
          setSecrets(data || []);
        }
      } catch (error) {
        console.error('Failed to fetch secrets', error);
      } finally {
        setLoading(false);
      }
    }

    fetchSecrets();
  }, [projectId, token]);

  async function handleUnlock(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setUnlocking(true);
    try {
      // 1. Fetch Vault & Project info
      const vcResponse = await fetch('/api/v1/vault/config', {
        headers: { Authorization: `Bearer ${token}` },
      });
      const vc = await vcResponse.json();

      const pResponse = await fetch(`/api/v1/projects/${projectId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      const project = await pResponse.json();

      // 2. Unwrap Master Key
      const salt = hexToBytes(vc.master_key_salt);
      const wrappedMK = hexToBytes(vc.wrapped_master_key);
      const adminKEK = await deriveKey(adminPassword, salt);
      const masterKey = await decrypt(adminKEK, wrappedMK);

      // 3. Unwrap Project Data Key
      const wrappedDK = hexToBytes(project.wrapped_data_key);
      const dataKey = await decrypt(masterKey, wrappedDK);

      setProjectDataKey(dataKey);
      setUnlocked(true);
      toast({
        title: 'Project unlocked',
        description: 'Secrets decrypted and ready for use.',
        color: 'teal',
      });

      // Decrypt all existing secrets
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
    } catch {
      toast({
        title: 'Decryption failed',
        description: 'Invalid admin password or corrupted key data.',
        color: 'red',
      });
    } finally {
      setUnlocking(false);
    }
  }

  async function handleAddSecret(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!projectDataKey) return;

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
          title: 'Secret saved',
          description: `Key ${newSecret.key} has been encrypted and stored.`,
          color: 'teal',
        });
        // Refresh secrets
        const refreshResponse = await fetch(
          `/api/v1/secrets?project_id=${projectId}`,
          {
            headers: { Authorization: `Bearer ${token}` },
          }
        );
        if (refreshResponse.ok) {
          const data = await refreshResponse.json();
          setSecrets(data || []);

          // Re-decrypt all
          const newDecrypted: Record<string, string> = {};
          for (const s of data) {
            try {
              const ciphertext = hexToBytes(s.value);
              const plaintext = await decrypt(projectDataKey, ciphertext);
              newDecrypted[s.id] = new TextDecoder().decode(plaintext);
            } catch {
              console.error('Failed to decrypt secret', s.id);
            }
          }
          setDecryptedValues(newDecrypted);
        }
      } else {
        throw new Error('Failed to save secret on server');
      }
    } catch (error) {
      console.error('Failed to add secret', error);
      toast({
        title: 'Save failed',
        description:
          error instanceof Error ? error.message : 'Encryption error',
        color: 'red',
      });
    } finally {
      setCreating(false);
    }
  }

  return (
    <Stack gap="6">
      <Button variant="text" onClick={() => navigate(-1)} p="0">
        <Flex gap="2" align="center">
          <IconChevronLeft size={18} />
          <Text>Back</Text>
        </Flex>
      </Button>

      {unlocked ? (
        <>
          <Flex justify="between" align="center">
            <Box>
              <Text size="6" weight="bold">
                Secrets
              </Text>
              <Text color="muted">
                End-to-end encrypted secrets for this environment.
              </Text>
            </Box>
            <Button variant="filled" onClick={() => setIsModalOpen(true)}>
              <IconPlus size={18} />
              <Text>Add Secret</Text>
            </Button>
          </Flex>

          <Card p="0" style={{ overflow: 'hidden' }}>
            <Table.Root>
              <Table.Header>
                <Table.Row>
                  <Table.ColumnHeader>Key</Table.ColumnHeader>
                  <Table.ColumnHeader>Value (Decrypted)</Table.ColumnHeader>
                  <Table.ColumnHeader>Version</Table.ColumnHeader>
                </Table.Row>
              </Table.Header>
              <Table.Body>
                {loading ? (
                  Array.from({ length: 3 }).map((_, i) => (
                    <Table.Row key={`loading-${i}`}>
                      <Table.Cell colSpan={3}>
                        <Box
                          p="4"
                          className="animate-pulse bg-surface-container-highest"
                          style={{ borderRadius: 'var(--pittorica-radius-sm)' }}
                        />
                      </Table.Cell>
                    </Table.Row>
                  ))
                ) : secrets.length === 0 ? (
                  <Table.Row>
                    <Table.Cell colSpan={3}>
                      <Flex p="8" justify="center">
                        <Text color="muted">No secrets found.</Text>
                      </Flex>
                    </Table.Cell>
                  </Table.Row>
                ) : (
                  secrets.map((s) => (
                    <Table.Row key={s.id}>
                      <Table.Cell>
                        <Text
                          weight="bold"
                          style={{ fontFamily: 'var(--pittorica-font-code)' }}
                        >
                          {s.key}
                        </Text>
                      </Table.Cell>
                      <Table.Cell>
                        <Box
                          p="2"
                          style={{
                            backgroundColor:
                              'var(--pittorica-color-surface-container)',
                            borderRadius: 'var(--pittorica-radius-sm)',
                          }}
                        >
                          <Text
                            style={{ fontFamily: 'var(--pittorica-font-code)' }}
                          >
                            {decryptedValues[s.id] || '••••••••'}
                          </Text>
                        </Box>
                      </Table.Cell>
                      <Table.Cell>
                        <Text size="1" color="muted">
                          v{s.version}
                        </Text>
                      </Table.Cell>
                    </Table.Row>
                  ))
                )}
              </Table.Body>
            </Table.Root>
          </Card>
        </>
      ) : (
        <Card p="8" style={{ maxWidth: '500px', margin: '0 auto' }}>
          <form onSubmit={handleUnlock}>
            <Stack gap="6" style={{ textAlign: 'center' }}>
              <Box color="cyan">
                <IconKey size={48} />
              </Box>
              <Box>
                <Text size="5" weight="bold">
                  Unlock Project
                </Text>
                <Text color="muted">
                  Enter your admin password to unwrap the data key and decrypt
                  secrets.
                </Text>
              </Box>
              <TextField.Root>
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
              <Button
                type="submit"
                variant="filled"
                disabled={unlocking}
                style={{ width: '100%' }}
              >
                {unlocking ? 'Decrypting...' : 'Decrypt Secrets'}
              </Button>
            </Stack>
          </form>
        </Card>
      )}

      <Dialog
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Add New Secret"
      >
        <form onSubmit={handleAddSecret}>
          <Stack gap="4">
            <TextField.Root>
              <TextField.Input
                placeholder="SECRET_KEY"
                value={newSecret.key}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewSecret({ ...newSecret, key: e.target.value })
                }
                required
              />
            </TextField.Root>
            <TextField.Root>
              <TextField.Input
                placeholder="Secret Value"
                value={newSecret.value}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setNewSecret({ ...newSecret, value: e.target.value })
                }
                required
              />
            </TextField.Root>
            <Flex justify="end" gap="3">
              <Button variant="text" onClick={() => setIsModalOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" variant="filled" disabled={creating}>
                {creating ? 'Encrypting...' : 'Encrypt & Save'}
              </Button>
            </Flex>
          </Stack>
        </form>
      </Dialog>
    </Stack>
  );
}
