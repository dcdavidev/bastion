import { useState, useEffect } from "react";
import { Box, Text, Stack, Card, Grid, Table } from "@pittorica/react";
import { useAuth } from "../../contexts/auth-context";
import { Shield, Users, Briefcase, Activity } from "lucide-react";

export default function Overview() {
  const [stats, setStats] = useState({ clients: 0, logs: 0 });
  const [latestLogs, setLatestLogs] = useState<any[]>([]);
  const { token, isAdmin } = useAuth();

  useEffect(() => {
    fetchStats();
  }, [token]);

  async function fetchStats() {
    try {
      const cResp = await fetch("/api/v1/clients", {
        headers: { "Authorization": `Bearer ${token}` }
      });
      const clients = await cResp.json();

      const lResp = await fetch("/api/v1/audit?limit=5", {
        headers: { "Authorization": `Bearer ${token}` }
      });
      const logs = await lResp.json();

      setStats({
        clients: clients?.length || 0,
        logs: logs?.length || 0,
      });
      setLatestLogs(logs || []);
    } catch (err) {
      console.error("Failed to fetch dashboard stats", err);
    }
  }

  return (
    <Stack gap="8">
      <Box>
        <Text size="6" weight="bold">Vault Overview</Text>
        <Text color="muted">Real-time status of your secure storage.</Text>
      </Box>

      <Grid columns="3" gap="6">
        <StatCard label="Total Clients" value={stats.clients.toString()} icon={<Users size={20} />} />
        <StatCard label="Recent Activities" value={stats.logs.toString()} icon={<Activity size={20} />} />
        <StatCard label="Security Level" value="Max" icon={<Shield size={20} />} />
      </Grid>

      {isAdmin && (
        <Stack gap="4">
          <Text size="4" weight="bold">Recent Audit Logs</Text>
          <Card padding="0" overflow="hidden">
            <Table>
              <Table.Header>
                <Table.Row>
                  <Table.HeaderCell>Time</Table.HeaderCell>
                  <Table.HeaderCell>Action</Table.HeaderCell>
                  <Table.HeaderCell>Target</Table.HeaderCell>
                </Table.Row>
              </Table.Header>
              <Table.Body>
                {latestLogs.map((log) => (
                  <Table.Row key={log.id}>
                    <Table.Cell><Text size="1" color="muted">{new Date(log.created_at).toLocaleTimeString()}</Text></Table.Cell>
                    <Table.Cell><Text weight="bold" size="2">{log.action}</Text></Table.Cell>
                    <Table.Cell><Text size="2" color="muted">{log.target_type}</Text></Table.Cell>
                  </Table.Row>
                ))}
              </Table.Body>
            </Table>
          </Card>
        </Stack>
      )}
    </Stack>
  );
}

function StatCard({ label, value, icon }: { label: string; value: string; icon: React.ReactNode }) {
  return (
    <Card padding="6">
      <Stack gap="3">
        <Box color="cyan" display="flex" alignItems="center" gap="2">
          {icon}
          <Text size="1" color="muted" weight="bold" transform="uppercase">{label}</Text>
        </Box>
        <Text size="8" weight="bold">{value}</Text>
      </Stack>
    </Card>
  );
}
