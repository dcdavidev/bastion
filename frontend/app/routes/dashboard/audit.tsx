import { useState, useEffect } from "react";
import { Box, Text, Stack, Card, Button, Table, Input } from "@pittorica/react";
import { useAuth } from "../../contexts/auth-context";
import { Filter, Calendar, Terminal } from "lucide-react";

interface AuditLog {
  id: string;
  action: string;
  target_type: string;
  target_id: string;
  metadata: Record<string, any>;
  created_at: string;
}

export default function AuditLogs() {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionFilter, setActionFilter] = useState("");
  const [targetFilter, setTargetFilter] = useState("");
  const { token } = useAuth();

  useEffect(() => {
    fetchLogs();
  }, [token, actionFilter, targetFilter]);

  async function fetchLogs() {
    setLoading(true);
    let url = `/api/v1/audit?limit=100`;
    if (actionFilter) url += `&action=${actionFilter}`;
    if (targetFilter) url += `&target_type=${targetFilter}`;

    try {
      const response = await fetch(url, {
        headers: { "Authorization": `Bearer ${token}` }
      });
      if (response.ok) {
        const data = await response.json();
        setLogs(data || []);
      }
    } catch (err) {
      console.error("Failed to fetch logs", err);
    } finally {
      setLoading(false);
    }
  }

  return (
    <Stack gap="6">
      <Box>
        <Text size="6" weight="bold">Audit Logs</Text>
        <Text color="muted">Full history of sensitive operations within the vault.</Text>
      </Box>

      <Card padding="4">
        <Stack direction="row" gap="4" alignItems="center">
          <Box display="flex" alignItems="center" gap="2" flex="1">
            <Filter size={16} color="muted" />
            <Input 
              placeholder="Filter by Action" 
              value={actionFilter}
              onChange={(e) => setActionFilter(e.target.value)}
            />
          </Box>
          <Box flex="1">
            <Input 
              placeholder="Filter by Target" 
              value={targetFilter}
              onChange={(e) => setTargetFilter(e.target.value)}
            />
          </Box>
          <Button variant="ghost" onClick={() => { setActionFilter(""); setTargetFilter(""); }}>Clear</Button>
        </Stack>
      </Card>

      <Card padding="0" overflow="hidden">
        <Table>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Timestamp</Table.HeaderCell>
              <Table.HeaderCell>Action</Table.HeaderCell>
              <Table.HeaderCell>Target</Table.HeaderCell>
              <Table.HeaderCell>Metadata</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {loading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <Table.Row key={i}>
                  <Table.Cell colSpan={4}><Box padding="4" className="animate-pulse bg-surface-container-highest rounded" /></Table.Cell>
                </Table.Row>
              ))
            ) : logs.length === 0 ? (
              <Table.Row><Table.Cell colSpan={4} textAlign="center"><Box padding="8"><Text color="muted">No logs found.</Text></Box></Table.Cell></Table.Row>
            ) : logs.map((log) => (
              <Table.Row key={log.id}>
                <Table.Cell>
                  <Stack direction="row" gap="2" alignItems="center">
                    <Calendar size={14} color="muted" />
                    <Text size="1" color="muted">{new Date(log.created_at).toLocaleString()}</Text>
                  </Stack>
                </Table.Cell>
                <Table.Cell>
                  <Box background="accent/10" paddingX="2" paddingY="1" borderRadius="sm" display="inline-block">
                    <Text weight="bold" size="1" color="cyan">{log.action}</Text>
                  </Box>
                </Table.Cell>
                <Table.Cell>
                  <Stack gap="0">
                    <Text size="2" weight="medium">{log.target_type}</Text>
                    <Text size="1" color="muted" family="mono">{log.target_id?.substring(0, 8)}...</Text>
                  </Stack>
                </Table.Cell>
                <Table.Cell>
                  <Box background="surface-container" padding="2" borderRadius="sm" maxHeight="100px" overflow="auto">
                    {log.metadata && Object.keys(log.metadata).length > 0 ? (
                      <Stack gap="1">
                        {Object.entries(log.metadata).map(([k, v]) => (
                          <Box key={k} display="flex" gap="2">
                            <Text size="1" weight="bold" color="muted">{k}:</Text>
                            <Text size="1" family="mono">{String(v)}</Text>
                          </Box>
                        ))}
                      </Stack>
                    ) : (
                      <Text size="1" color="muted">No metadata</Text>
                    )}
                  </Box>
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      </Card>
    </Stack>
  );
}
