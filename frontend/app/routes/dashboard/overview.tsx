import { useEffect, useState } from 'react';

import { IconActivity, IconShield, IconUsers } from '@tabler/icons-react';

import { Box, Card, Flex, Grid, Stack, Table, Text } from '@pittorica/react';

import { useAuth } from '../../contexts/auth-context';

interface AuditLog {
  id: string;
  action: string;
  target_type: string;
  target_id: string;
  created_at: string;
}

export default function Overview() {
  const [stats, setStats] = useState({ clients: 0, logs: 0 });
  const [latestLogs, setLatestLogs] = useState<AuditLog[]>([]);
  const { token, isAdmin } = useAuth();

  useEffect(() => {
    async function fetchStats() {
      try {
        const cResp = await fetch('/api/v1/clients', {
          headers: { Authorization: `Bearer ${token}` },
        });
        const clients = await cResp.json();

        const lResp = await fetch('/api/v1/audit?limit=5', {
          headers: { Authorization: `Bearer ${token}` },
        });
        const logs = await lResp.json();

        setStats({
          clients: clients?.length || 0,
          logs: logs?.length || 0,
        });
        setLatestLogs(logs || []);
      } catch (error) {
        console.error('Failed to fetch dashboard stats', error);
      }
    }

    fetchStats();
  }, [token]);

  return (
    <Stack gap="8">
      <Box>
        <Text size="6" weight="bold">
          Vault Overview
        </Text>
        <Text color="muted">Real-time status of your secure storage.</Text>
      </Box>

      <Grid columns="3" gap="6">
        <StatCard
          label="Total Clients"
          value={stats.clients.toString()}
          icon={<IconUsers size={20} />}
        />
        <StatCard
          label="Recent Activities"
          value={stats.logs.toString()}
          icon={<IconActivity size={20} />}
        />
        <StatCard
          label="Security Level"
          value="Max"
          icon={<IconShield size={20} />}
        />
      </Grid>

      {isAdmin && (
        <Stack gap="4">
          <Text size="4" weight="bold">
            Recent Audit Logs
          </Text>
          <Card p="0" style={{ overflow: 'hidden' }}>
            <Table.Root>
              <Table.Header>
                <Table.Row>
                  <Table.ColumnHeader>Time</Table.ColumnHeader>
                  <Table.ColumnHeader>Action</Table.ColumnHeader>
                  <Table.ColumnHeader>Target</Table.ColumnHeader>
                </Table.Row>
              </Table.Header>
              <Table.Body>
                {latestLogs.map((log) => (
                  <Table.Row key={log.id}>
                    <Table.Cell>
                      <Text size="1" color="muted">
                        {new Date(log.created_at).toLocaleTimeString()}
                      </Text>
                    </Table.Cell>
                    <Table.Cell>
                      <Text weight="bold" size="2">
                        {log.action}
                      </Text>
                    </Table.Cell>
                    <Table.Cell>
                      <Text size="2" color="muted">
                        {log.target_type}
                      </Text>
                    </Table.Cell>
                  </Table.Row>
                ))}
              </Table.Body>
            </Table.Root>
          </Card>
        </Stack>
      )}
    </Stack>
  );
}

function StatCard({
  label,
  value,
  icon,
}: {
  label: string;
  value: string;
  icon: React.ReactNode;
}) {
  return (
    <Card p="6">
      <Stack gap="3">
        <Flex color="cyan" align="center" gap="2">
          {icon}
          <Text
            size="1"
            color="muted"
            weight="bold"
            style={{ textTransform: 'uppercase' }}
          >
            {label}
          </Text>
        </Flex>
        <Text size="8" weight="bold">
          {value}
        </Text>
      </Stack>
    </Card>
  );
}
