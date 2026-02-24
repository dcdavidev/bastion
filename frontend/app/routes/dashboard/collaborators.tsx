import { useState, useEffect } from "react";
import { Box, Text, Stack, Card, Button, Table, Dialog, Input } from "@pittorica/react";
import { useAuth } from "../../contexts/auth-context";
import { Plus, UserPlus, Shield } from "lucide-react";

export default function Collaborators() {
  const [loading, setLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newCollab, setNewCollab] = useState({ username: "", password: "", projectId: "" });
  const [creating, setCreating] = useState(false);
  const { token } = useAuth();

  // Note: We'd typically fetch a list of collaborators here.
  // For the sake of this step, we'll focus on the creation UI.

  async function handleCreateCollaborator(e: React.FormEvent) {
    e.preventDefault();
    setCreating(true);
    
    // In a real E2EE implementation, the admin would:
    // 1. Fetch the project data key (already done in projects view logic)
    // 2. Wrap it with the collaborator's password (derived KEK)
    const mockWrappedKey = "collab-wrapped-key-" + Math.random().toString(16);

    try {
      const response = await fetch("/api/v1/collaborators", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({
          username: newCollab.username,
          password_hash: "mock-hash", // Placeholder
          salt: "mock-salt",         // Placeholder
          project_id: newCollab.projectId,
          wrapped_data_key: mockWrappedKey
        }),
      });

      if (response.ok) {
        setIsModalOpen(false);
        setNewCollab({ username: "", password: "", projectId: "" });
      }
    } catch (err) {
      console.error("Failed to create collaborator", err);
    } finally {
      setCreating(false);
    }
  }

  return (
    <Stack gap="6">
      <Box display="flex" justifyContent="space-between" alignItems="center">
        <Box>
          <Text size="6" weight="bold">Collaborators</Text>
          <Text color="muted">Manage restricted access for team members.</Text>
        </Box>
        <Button variant="primary" onClick={() => setIsModalOpen(true)}>
          <Stack direction="row" gap="2" alignItems="center">
            <UserPlus size={18} />
            <Text>New Collaborator</Text>
          </Stack>
        </Button>
      </Box>

      <Card padding="8" textAlign="center">
        <Stack gap="4" alignItems="center">
          <Box color="cyan"><Shield size={48} /></Box>
          <Box>
            <Text size="4" weight="bold">Access Control</Text>
            <Text color="muted">Collaborators can only see and use secrets for projects they are assigned to.</Text>
          </Box>
        </Stack>
      </Card>

      <Dialog open={isModalOpen} onClose={() => setIsModalOpen(false)} title="Add Collaborator">
        <form onSubmit={handleCreateCollaborator}>
          <Stack gap="4">
            <Input 
              placeholder="Username" 
              value={newCollab.username} 
              onChange={(e) => setNewCollab({...newCollab, username: e.target.value})} 
              required 
            />
            <Input 
              type="password" 
              placeholder="Assign a Password" 
              value={newCollab.password} 
              onChange={(e) => setNewCollab({...newCollab, password: e.target.value})} 
              required 
            />
            <Input 
              placeholder="Project ID (UUID)" 
              value={newCollab.projectId} 
              onChange={(e) => setNewCollab({...newCollab, projectId: e.target.value})} 
              required 
            />
            <Box display="flex" justifyContent="flex-end" gap="3">
              <Button variant="ghost" onClick={() => setIsModalOpen(false)}>Cancel</Button>
              <Button type="submit" variant="primary" loading={creating}>Grant Access</Button>
            </Box>
          </Stack>
        </form>
      </Dialog>
    </Stack>
  );
}
