import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router";
import { Box, Text, Stack, Card, Button, Table, Dialog, Input } from "@pittorica/react";
import { useAuth } from "../../contexts/auth-context";
import { Plus, Briefcase, ChevronLeft } from "lucide-react";
import { deriveKey, encrypt, decrypt, hexToBytes, bytesToHex } from "../../utils/crypto";

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
  const [newProjectName, setNewProjectName] = useState("");
  const [adminPassword, setAdminPassword] = useState("");
  const [creating, setCreating] = useState(false);
  const { token } = useAuth();

  useEffect(() => {
    fetchProjects();
  }, [clientId, token]);

  async function fetchProjects() {
    try {
      const response = await fetch(`/api/v1/projects?client_id=${clientId}`, {
        headers: { "Authorization": `Bearer ${token}` }
      });
      if (response.ok) {
        const data = await response.json();
        setProjects(data || []);
      }
    } catch (err) {
      console.error("Failed to fetch projects", err);
    } finally {
      setLoading(false);
    }
  }

  async function handleCreateProject(e: React.FormEvent) {
    e.preventDefault();
    setCreating(true);

    try {
      // 1. Fetch Vault Config
      const vcResponse = await fetch("/api/v1/vault/config", {
        headers: { "Authorization": `Bearer ${token}` }
      });
      const vc = await vcResponse.json();

      // 2. Unwrap Master Key
      const salt = hexToBytes(vc.master_key_salt);
      const wrappedMK = hexToBytes(vc.wrapped_master_key);
      const adminKEK = await deriveKey(adminPassword, salt);
      const masterKey = await decrypt(adminKEK, wrappedMK);

      // 3. Generate New Project Data Key
      const dataKey = window.crypto.getRandomValues(new Uint8Array(32));
      const wrappedDK = await encrypt(masterKey, dataKey);

      // 4. Send to server
      const response = await fetch("/api/v1/projects", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({ 
          client_id: clientId, 
          name: newProjectName,
          wrapped_data_key: bytesToHex(wrappedDK)
        }),
      });

      if (response.ok) {
        setNewProjectName("");
        setAdminPassword("");
        setIsModalOpen(false);
        fetchProjects();
      } else {
        throw new Error("Failed to create project on server");
      }
    } catch (err: any) {
      console.error("Encryption or request failed", err);
      alert(err.message || "Invalid admin password or crypto error");
    } finally {
      setCreating(false);
    }
  }

  return (
    <Stack gap="6">
      <Button variant="ghost" onClick={() => navigate("/dashboard/clients")} padding="0">
        <Stack direction="row" gap="2" alignItems="center">
          <ChevronLeft size={18} />
          <Text>Back to Clients</Text>
        </Stack>
      </Button>

      <Box display="flex" justifyContent="space-between" alignItems="center">
        <Box>
          <Text size="6" weight="bold">Projects</Text>
          <Text color="muted">Manage isolated environments and secrets for this client.</Text>
        </Box>
        <Button variant="primary" onClick={() => setIsModalOpen(true)}>
          <Stack direction="row" gap="2" alignItems="center">
            <Plus size={18} />
            <Text>New Project</Text>
          </Stack>
        </Button>
      </Box>

      <Card padding="0" overflow="hidden">
        <Table>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Project Name</Table.HeaderCell>
              <Table.HeaderCell>ID</Table.HeaderCell>
              <Table.HeaderCell>Created At</Table.HeaderCell>
              <Table.HeaderCell textAlign="right">Actions</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {loading ? (
              <Table.Row><Table.Cell colSpan={4} textAlign="center"><Box padding="8"><Text color="muted">Loading projects...</Text></Box></Table.Cell></Table.Row>
            ) : projects.length === 0 ? (
              <Table.Row><Table.Cell colSpan={4} textAlign="center"><Box padding="8"><Text color="muted">No projects found.</Text></Box></Table.Cell></Table.Row>
            ) : projects.map((p) => (
              <Table.Row key={p.id}>
                <Table.Cell>
                  <Stack direction="row" gap="3" alignItems="center">
                    <Box color="cyan"><Briefcase size={16} /></Box>
                    <Text weight="medium">{p.name}</Text>
                  </Stack>
                </Table.Cell>
                <Table.Cell><Text family="mono" size="1" color="muted">{p.id}</Text></Table.Cell>
                <Table.Cell><Text size="2">{new Date(p.created_at).toLocaleDateString()}</Text></Table.Cell>
                <Table.Cell textAlign="right">
                  <Button variant="ghost" size="small" onClick={() => navigate(`/dashboard/projects/${p.id}`)}>
                    View Secrets
                  </Button>
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      </Card>

      <Dialog open={isModalOpen} onClose={() => setIsModalOpen(false)} title="Create New Project">
        <form onSubmit={handleCreateProject}>
          <Stack gap="4">
            <Text color="muted">This will generate a new unique E2EE data key for the project.</Text>
            <Input placeholder="Project Name" value={newProjectName} onChange={(e) => setNewProjectName(e.target.value)} required />
            <Input type="password" placeholder="Admin Password (to wrap key)" value={adminPassword} onChange={(e) => setAdminPassword(e.target.value)} required />
            <Box display="flex" justifyContent="flex-end" gap="3">
              <Button variant="ghost" onClick={() => setIsModalOpen(false)}>Cancel</Button>
              <Button type="submit" variant="primary" loading={creating}>Create & Secure</Button>
            </Box>
          </Stack>
        </form>
      </Dialog>
    </Stack>
  );
}
