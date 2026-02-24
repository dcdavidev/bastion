import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router";
import { Box, Text, Stack, Card, Button, Table, Dialog, Input, Label } from "@pittorica/react";
import { useAuth } from "../../contexts/auth-context";
import { Plus, Key, ChevronLeft, Eye, EyeOff } from "lucide-react";
import { deriveKey, encrypt, decrypt, hexToBytes, bytesToHex } from "../../utils/crypto";

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
  const [decryptedValues, setDecryptedValues] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newSecret, setNewSecret] = useState({ key: "", value: "" });
  const [adminPassword, setAdminPassword] = useState("");
  const [unlocking, setUnlocking] = useState(false);
  const [unlocked, setUnlocked] = useState(false);
  const [projectDataKey, setProjectDataKey] = useState<Uint8Array | null>(null);
  const { token } = useAuth();

  useEffect(() => {
    fetchSecrets();
  }, [projectId, token]);

  async function fetchSecrets() {
    try {
      const response = await fetch(`/api/v1/secrets?project_id=${projectId}`, {
        headers: { "Authorization": `Bearer ${token}` }
      });
      if (response.ok) {
        const data = await response.json();
        setSecrets(data || []);
      }
    } catch (err) {
      console.error("Failed to fetch secrets", err);
    } finally {
      setLoading(false);
    }
  }

  async function handleUnlock(e: React.FormEvent) {
    e.preventDefault();
    setUnlocking(true);
    try {
      // 1. Fetch Vault & Project info
      const vcResponse = await fetch("/api/v1/vault/config", {
        headers: { "Authorization": `Bearer ${token}` }
      });
      const vc = await vcResponse.json();

      const pResponse = await fetch(`/api/v1/projects/${projectId}`, {
        headers: { "Authorization": `Bearer ${token}` }
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
      
      // Decrypt all existing secrets
      const newDecrypted: Record<string, string> = {};
      for (const s of secrets) {
        const ciphertext = hexToBytes(s.value);
        const plaintext = await decrypt(dataKey, ciphertext);
        newDecrypted[s.id] = new TextDecoder().decode(plaintext);
      }
      setDecryptedValues(newDecrypted);
    } catch (err) {
      alert("Invalid password or decryption error");
    } finally {
      setUnlocking(false);
    }
  }

  async function handleAddSecret(e: React.FormEvent) {
    e.preventDefault();
    if (!projectDataKey) return;

    try {
      const plaintext = new TextEncoder().encode(newSecret.value);
      const ciphertext = await encrypt(projectDataKey, plaintext);

      const response = await fetch("/api/v1/secrets", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
          project_id: projectId,
          key: newSecret.key,
          value: bytesToHex(ciphertext)
        }),
      });

      if (response.ok) {
        setNewSecret({ key: "", value: "" });
        setIsModalOpen(false);
        fetchSecrets();
      }
    } catch (err) {
      console.error("Failed to add secret", err);
    }
  }

  return (
    <Stack gap="6">
      <Button variant="ghost" onClick={() => navigate(-1)} padding="0">
        <Stack direction="row" gap="2" alignItems="center">
          <ChevronLeft size={18} />
          <Text>Back</Text>
        </Stack>
      </Button>

      {!unlocked ? (
        <Card padding="8" maxWidth="500px" margin="0 auto">
          <form onSubmit={handleUnlock}>
            <Stack gap="6" textAlign="center">
              <Box color="cyan"><Key size={48} /></Box>
              <Box>
                <Text size="5" weight="bold">Unlock Project</Text>
                <Text color="muted">Enter your admin password to unwrap the data key and decrypt secrets.</Text>
              </Box>
              <Input 
                type="password" 
                placeholder="Admin Password" 
                value={adminPassword} 
                onChange={(e) => setAdminPassword(e.target.value)} 
                required 
              />
              <Button type="submit" variant="primary" loading={unlocking} width="100%">
                Decrypt Secrets
              </Button>
            </Stack>
          </form>
        </Card>
      ) : (
        <>
          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Box>
              <Text size="6" weight="bold">Secrets</Text>
              <Text color="muted">End-to-end encrypted secrets for this environment.</Text>
            </Box>
            <Button variant="primary" onClick={() => setIsModalOpen(true)}>
              <Plus size={18} />
              <Text>Add Secret</Text>
            </Button>
          </Box>

          <Card padding="0" overflow="hidden">
            <Table>
              <Table.Header>
                <Table.Row>
                  <Table.HeaderCell>Key</Table.HeaderCell>
                  <Table.HeaderCell>Value (Decrypted)</Table.HeaderCell>
                  <Table.HeaderCell>Version</Table.HeaderCell>
                </Table.Row>
              </Table.Header>
              <Table.Body>
                {secrets.map((s) => (
                  <Table.Row key={s.id}>
                    <Table.Cell><Text weight="bold" family="mono">{s.key}</Text></Table.Cell>
                    <Table.Cell>
                      <Box background="surface-container" padding="2" borderRadius="sm">
                        <Text family="mono">{decryptedValues[s.id] || "••••••••"}</Text>
                      </Box>
                    </Table.Cell>
                    <Table.Cell><Text size="1" color="muted">v{s.version}</Text></Table.Cell>
                  </Table.Row>
                ))}
              </Table.Body>
            </Table>
          </Card>
        </>
      )}

      <Dialog open={isModalOpen} onClose={() => setIsModalOpen(false)} title="Add New Secret">
        <form onSubmit={handleAddSecret}>
          <Stack gap="4">
            <Input 
              placeholder="SECRET_KEY" 
              value={newSecret.key} 
              onChange={(e) => setNewSecret({...newSecret, key: e.target.value})} 
              required 
            />
            <Input 
              placeholder="Secret Value" 
              value={newSecret.value} 
              onChange={(e) => setNewSecret({...newSecret, value: e.target.value})} 
              required 
            />
            <Box display="flex" justifyContent="flex-end" gap="3">
              <Button variant="ghost" onClick={() => setIsModalOpen(false)}>Cancel</Button>
              <Button type="submit" variant="primary">Encrypt & Save</Button>
            </Box>
          </Stack>
        </form>
      </Dialog>
    </Stack>
  );
}
