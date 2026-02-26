import { useEffect, useState } from 'react';

import {
  IconActivity,
  IconHistory,
  IconShieldCheck,
  IconUsers,
} from '@tabler/icons-react';

import {
  Badge,
  Box,
  Card,
  Chip,
  Flex,
  Grid,
  Stack,
  Table,
  Text,
} from '@pittorica/react';

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
  const [loading, setLoading] = useState(true);
  const { token, isAdmin } = useAuth();

  useEffect(() => {
    async function fetchStats() {
      setLoading(true);
      try {
        const [cResp, lResp] = await Promise.all([
          fetch('/api/v1/clients', {
            headers: { Authorization: `Bearer ${token}` },
          }),
          fetch('/api/v1/audit?limit=5', {
            headers: { Authorization: `Bearer ${token}` },
          }),
        ]);

        const clients = await cResp.json();
        const logs = await lResp.json();

        setStats({
          clients: clients?.length || 0,
          logs: logs?.length || 0,
        });
        setLatestLogs(logs || []);
      } catch (error) {
        console.error('Failed to fetch dashboard stats', error);
      } finally {
        setLoading(false);
      }
    }

    fetchStats();
  }, [token]);

  return (
    <Stack gap="8">
      <Box>
        <Text size="7" weight="bold" color="source">
          Vault Overview
        </Text>
        <Text color="muted" size="2">
          Monitor the health and activity of your secure fortress.
        </Text>
      </Box>

      <Grid columns="3" gap="6">
        <StatCard
          label="Active Clients"
          value={loading ? '...' : stats.clients.toString()}
          icon={<IconUsers size={24} />}
          description="Total registered tenants"
        />
        <StatCard
          label="System Events"
          value={loading ? '...' : stats.logs.toString()}
          icon={<IconActivity size={24} />}
          description="Actions in the last 24h"
        />
        <StatCard
          label="Encryption"
          value="AES-256"
          icon={<IconShieldCheck size={24} />}
          description="E2EE Protection Active"
        />
      </Grid>

      {isAdmin && (
        <Stack gap="4">
          <Flex align="center" gap="2">
            <IconHistory size={20} color="var(--pittorica-color-source)" />
            <Text size="4" weight="bold">
              Recent Activity
            </Text>
          </Flex>
          <Card p="0" style={{ overflow: 'hidden' }}>
            <Table.Root>
              <Table.Header>
                <Table.Row>
                  <Table.ColumnHeader>Time</Table.ColumnHeader>
                  <Table.ColumnHeader>Action</Table.ColumnHeader>
                  <Table.ColumnHeader>Target Type</Table.ColumnHeader>
                  <Table.ColumnHeader>Status</Table.ColumnHeader>
                </Table.Row>
              </Table.Header>
              <Table.Body>
                {latestLogs.length === 0 && !loading ? (
                  <Table.Row>
                    <Table.Cell colSpan={4}>
                      <Flex p="6" justify="center">
                        <Text color="muted">No recent activity found.</Text>
                      </Flex>
                    </Table.Cell>
                  </Table.Row>
                ) : (
                  latestLogs.map((log) => (
                    <Table.Row key={log.id}>
                      <Table.Cell>
                        <Text size="1" color="muted">
                          {new Date(log.created_at).toLocaleTimeString([], {
                            hour: '2-digit',
                            minute: '2-digit',
                          })}
                        </Text>
                      </Table.Cell>
                      <Table.Cell>
                        <Text weight="bold" size="2">
                          {log.action}
                        </Text>
                      </Table.Cell>
                      <Table.Cell>
                        <Badge variant="standard">{log.target_type}</Badge>
                      </Table.Cell>
                      <Table.Cell>
                        <Badge color="teal" variant="standard">
                          Success
                        </Badge>
                      </Table.Cell>
                    </Table.Row>
                  ))
                )}
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
  description,
}: {
  label: string;
  value: string;
  icon: React.ReactNode;
  description: string;
}) {
  return (
    <Card p="6">
      <Stack gap="4">
        <Flex justify="between" align="start">
          <Box color="source" p="2">
            {icon}
          </Box>
          <Chip variant="soft">LIVE</Chip>
        </Flex>
        <Stack gap="1">
          <Text size="8" weight="bold" style={{ lineHeight: 1 }}>
            {value}
          </Text>
          <Text size="2" weight="medium">
            {label}
          </Text>
          <Text size="1" color="muted">
            {description}
          </Text>
        </Stack>
      </Stack>
    </Card>
  );
}
