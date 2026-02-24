import { useState, useEffect } from "react";
import { Box, Text, Stack, Card, Button, Table } from "@pittorica/react";
import { useAuth } from "../../contexts/auth-context";
import { Plus, User } from "lucide-react";

interface Client {
  id: string;
  name: string;
  created_at: string;
}

export default function Clients() {
  const [clients, setClients] = useState<Client[]>([]);
  const [loading, setLoading] = useState(true);
  const { token } = useAuth();

  useEffect(() => {
    fetchClients();
  }, [token]);

  async function fetchClients() {
    try {
      const response = await fetch("/api/v1/clients", {
        headers: {
          "Authorization": `Bearer ${token}`
        }
      });
      if (response.ok) {
        const data = await response.json();
        setClients(data || []);
      }
    } catch (err) {
      console.error("Failed to fetch clients", err);
    } finally {
      setLoading(false);
    }
  }

  return (
    <Stack gap="6">
      <Box display="flex" justifyContent="space-between" alignItems="center">
        <Box>
          <Text size="6" weight="bold">Clients</Text>
          <Text color="muted">Manage your client base and their associated projects.</Text>
        </Box>
        <Button variant="primary">
          <Stack direction="row" gap="2" alignItems="center">
            <Plus size={18} />
            <Text>New Client</Text>
          </Stack>
        </Button>
      </Box>

      <Card padding="0" overflow="hidden">
        <Table>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Client Name</Table.HeaderCell>
              <Table.HeaderCell>ID</Table.HeaderCell>
              <Table.HeaderCell>Created At</Table.HeaderCell>
              <Table.HeaderCell textAlign="right">Actions</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {loading ? (
              <Table.Row>
                <Table.Cell colSpan={4}>
                  <Box padding="8" textAlign="center">
                    <Text color="muted">Loading clients...</Text>
                  </Box>
                </Table.Cell>
              </Table.Row>
            ) : clients.length === 0 ? (
              <Table.Row>
                <Table.Cell colSpan={4}>
                  <Box padding="8" textAlign="center">
                    <Text color="muted">No clients found.</Text>
                  </Box>
                </Table.Cell>
              </Table.Row>
            ) : clients.map((client) => (
              <Table.Row key={client.id}>
                <Table.Cell>
                  <Stack direction="row" gap="3" alignItems="center">
                    <Box 
                      width="32px" 
                      height="32px" 
                      borderRadius="full" 
                      background="accent/10" 
                      display="flex" 
                      alignItems="center" 
                      justifyContent="center"
                      color="cyan"
                    >
                      <User size={16} />
                    </Box>
                    <Text weight="medium">{client.name}</Text>
                  </Stack>
                </Table.Cell>
                <Table.Cell>
                  <Text family="mono" size="1" color="muted">{client.id}</Text>
                </Table.Cell>
                <Table.Cell>
                  <Text size="2">{new Date(client.created_at).toLocaleDateString()}</Text>
                </Table.Cell>
                <Table.Cell textAlign="right">
                  <Button variant="ghost" size="small">View Projects</Button>
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      </Card>
    </Stack>
  );
}
