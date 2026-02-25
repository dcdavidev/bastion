import type { ChangeEvent } from 'react';
import { useEffect, useState } from 'react';

import { IconCalendar, IconFilter } from '@tabler/icons-react';

import {
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
        <Text size="6" weight="bold">
          Audit Logs
        </Text>
        <Text color="muted">
          Full history of sensitive operations within the vault.
        </Text>
      </Box>

      <Card p="4">
        <Flex gap="4" align="center">
          <Flex align="center" gap="2" style={{ flex: 1 }}>
            <IconFilter size={16} color="var(--pittorica-color-muted)" />
            <TextField.Root style={{ flex: 1 }}>
              <TextField.Input
                placeholder="Filter by Action"
                value={actionFilter}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setActionFilter(e.target.value)
                }
              />
            </TextField.Root>
          </Flex>
          <Box style={{ flex: 1 }}>
            <TextField.Root>
              <TextField.Input
                placeholder="Filter by Target"
                value={targetFilter}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setTargetFilter(e.target.value)
                }
              />
            </TextField.Root>
          </Box>
          <Button
            variant="text"
            onClick={() => {
              setActionFilter('');
              setTargetFilter('');
            }}
          >
            Clear
          </Button>
        </Flex>
      </Card>

      <Card p="0" style={{ overflow: 'hidden' }}>
        <Table.Root>
          <Table.Header>
            <Table.Row>
              <Table.ColumnHeader>Timestamp</Table.ColumnHeader>
              <Table.ColumnHeader>Action</Table.ColumnHeader>
              <Table.ColumnHeader>Target</Table.ColumnHeader>
              <Table.ColumnHeader>Metadata</Table.ColumnHeader>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {loading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <Table.Row key={`loading-${i}`}>
                  <Table.Cell colSpan={4}>
                    <Box
                      p="4"
                      className="animate-pulse bg-surface-container-highest"
                      style={{ borderRadius: 'var(--pittorica-radius-sm)' }}
                    />
                  </Table.Cell>
                </Table.Row>
              ))
            ) : logs.length === 0 ? (
              <Table.Row>
                <Table.Cell colSpan={4}>
                  <Flex p="8" justify="center">
                    <Text color="muted">No logs found.</Text>
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
                      <Text size="1" color="muted">
                        {new Date(log.created_at).toLocaleString()}
                      </Text>
                    </Flex>
                  </Table.Cell>
                  <Table.Cell>
                    <Box
                      p="1"
                      px="2"
                      style={{
                        backgroundColor:
                          'rgba(var(--pittorica-color-accent-rgb), 0.1)',
                        borderRadius: 'var(--pittorica-radius-sm)',
                      }}
                      display="inline-block"
                    >
                      <Text weight="bold" size="1" color="cyan">
                        {log.action}
                      </Text>
                    </Box>
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
                        {log.target_id?.slice(0, 8)}...
                      </Text>
                    </Stack>
                  </Table.Cell>
                  <Table.Cell>
                    <Box
                      p="2"
                      style={{
                        backgroundColor:
                          'var(--pittorica-color-surface-container)',
                        borderRadius: 'var(--pittorica-radius-sm)',
                        maxHeight: '100px',
                        overflow: 'auto',
                      }}
                    >
                      {log.metadata && Object.keys(log.metadata).length > 0 ? (
                        <Stack gap="1">
                          {Object.entries(log.metadata).map(([k, v]) => (
                            <Flex key={k} gap="2">
                              <Text size="1" weight="bold" color="muted">
                                {k}:
                              </Text>
                              <Text
                                size="1"
                                style={{
                                  fontFamily: 'var(--pittorica-font-code)',
                                }}
                              >
                                {String(v)}
                              </Text>
                            </Flex>
                          ))}
                        </Stack>
                      ) : (
                        <Text size="1" color="muted">
                          No metadata
                        </Text>
                      )}
                    </Box>
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
