import type { ChangeEvent } from 'react';
import { useEffect, useState } from 'react';

import {
  IconCalendar,
  IconDeviceDesktop,
  IconFilter,
  IconInfoCircle,
} from '@tabler/icons-react';

import {
  Badge,
  Box,
  Button,
  Card,
  Flex,
  Stack,
  Table,
  Text,
  TextField,
} from '@pittorica/react';

import { useAuth } from '../../contexts/auth-context';

interface AuditLog {
  id: string;
  action: string;
  target_type: string;
  target_id: string;
  metadata: Record<string, string | number | boolean | null>;
  created_at: string;
}

export default function AuditLogs() {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionFilter, setActionFilter] = useState('');
  const [targetFilter, setTargetFilter] = useState('');
  const { token } = useAuth();

  useEffect(() => {
    async function fetchLogs() {
      setLoading(true);
      let url = `/api/v1/audit?limit=100`;
      if (actionFilter) url += `&action=${encodeURIComponent(actionFilter)}`;
      if (targetFilter)
        url += `&target_type=${encodeURIComponent(targetFilter)}`;

      try {
        const response = await fetch(url, {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (response.ok) {
          const data = await response.json();
          setLogs(data || []);
        }
      } catch (error) {
        console.error('Failed to fetch logs', error);
      } finally {
        setLoading(false);
      }
    }

    fetchLogs();
  }, [token, actionFilter, targetFilter]);

  return (
    <Stack gap="6">
      <Box>
        <Text size="7" weight="bold" color="source">
          Audit Logs
        </Text>
        <Text color="muted" size="2">
          Historical record of all operations and security events.
        </Text>
      </Box>

      <Card p="4">
        <Flex gap="4" align="end">
          <Box style={{ flex: 1 }}>
            <TextField.Root size="md" label="Action Filter">
              <TextField.Slot>
                <IconFilter size={16} color="var(--pittorica-color-muted)" />
              </TextField.Slot>
              <TextField.Input
                placeholder="e.g. LOGIN, SECRET_READ"
                value={actionFilter}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setActionFilter(e.target.value)
                }
              />
            </TextField.Root>
          </Box>
          <Box style={{ flex: 1 }}>
            <TextField.Root size="md" label="Target Type">
              <TextField.Slot>
                <IconDeviceDesktop
                  size={16}
                  color="var(--pittorica-color-muted)"
                />
              </TextField.Slot>
              <TextField.Input
                placeholder="e.g. PROJECT, CLIENT"
                value={targetFilter}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setTargetFilter(e.target.value)
                }
              />
            </TextField.Root>
          </Box>
          <Button
            variant="text"
            size="md"
            onClick={() => {
              setActionFilter('');
              setTargetFilter('');
            }}
          >
            Reset Filters
          </Button>
        </Flex>
      </Card>

      <Card p="0" style={{ overflow: 'hidden' }}>
        <Table.Root>
          <Table.Header>
            <Table.Row>
              <Table.ColumnHeader>Timestamp</Table.ColumnHeader>
              <Table.ColumnHeader>Operation</Table.ColumnHeader>
              <Table.ColumnHeader>Target</Table.ColumnHeader>
              <Table.ColumnHeader>Context / Metadata</Table.ColumnHeader>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {loading ? (
              Array.from({ length: 5 }).map((_, i) => (
                // eslint-disable-next-line @eslint-react/no-array-index-key
                <Table.Row key={`loading-${i}`}>
                  <Table.Cell colSpan={4}>
                    <Box
                      p="4"
                      className="animate-pulse bg-surface-container"
                      style={{
                        borderRadius: 'var(--pittorica-radius-sm)',
                        height: '40px',
                      }}
                    />
                  </Table.Cell>
                </Table.Row>
              ))
            ) : logs.length === 0 ? (
              <Table.Row>
                <Table.Cell colSpan={4}>
                  <Flex p="8" justify="center">
                    <Stack align="center" gap="2">
                      <IconInfoCircle
                        size={32}
                        color="var(--pittorica-color-muted)"
                      />
                      <Text color="muted">
                        No events recorded for these filters.
                      </Text>
                    </Stack>
                  </Flex>
                </Table.Cell>
              </Table.Row>
            ) : (
              logs.map((log) => (
                <Table.Row key={log.id}>
                  <Table.Cell>
                    <Flex gap="2" align="center">
                      <IconCalendar
                        size={14}
                        color="var(--pittorica-color-muted)"
                      />
                      <Stack gap="0">
                        <Text size="1" weight="bold">
                          {new Date(log.created_at).toLocaleDateString()}
                        </Text>
                        <Text size="1" color="muted">
                          {new Date(log.created_at).toLocaleTimeString()}
                        </Text>
                      </Stack>
                    </Flex>
                  </Table.Cell>
                  <Table.Cell>
                    <Badge
                      variant="standard"
                      style={{ textTransform: 'uppercase' }}
                    >
                      {log.action}
                    </Badge>
                  </Table.Cell>
                  <Table.Cell>
                    <Stack gap="0">
                      <Text size="2" weight="medium">
                        {log.target_type}
                      </Text>
                      <Text
                        size="1"
                        color="muted"
                        style={{ fontFamily: 'var(--pittorica-font-code)' }}
                      >
                        {log.target_id?.slice(0, 13)}...
                      </Text>
                    </Stack>
                  </Table.Cell>
                  <Table.Cell>
                    {log.metadata && Object.keys(log.metadata).length > 0 ? (
                      <Flex gap="1" wrap="wrap">
                        {Object.entries(log.metadata).map(([k, v]) => (
                          <Badge key={k} variant="standard" color="muted">
                            {k}: {String(v)}
                          </Badge>
                        ))}
                      </Flex>
                    ) : (
                      <Text
                        size="1"
                        color="muted"
                        style={{ fontStyle: 'italic' }}
                      >
                        None
                      </Text>
                    )}
                  </Table.Cell>
                </Table.Row>
              ))
            )}
          </Table.Body>
        </Table.Root>
      </Card>
    </Stack>
  );
}
