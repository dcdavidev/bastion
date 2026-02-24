import { Box, Text, Stack, Card, Grid } from "@pittorica/react";

export default function Overview() {
  return (
    <Stack gap="6">
      <Box>
        <Text size="6" weight="bold">Vault Overview</Text>
        <Text color="muted">Welcome back to your secure fortress.</Text>
      </Box>

      <Grid columns="3" gap="6">
        <StatCard label="Total Clients" value="0" />
        <StatCard label="Active Projects" value="0" />
        <StatCard label="Secrets Managed" value="0" />
      </Grid>
    </Stack>
  );
}

function StatCard({ label, value }: { label: string; value: string }) {
  return (
    <Card padding="6">
      <Stack gap="1">
        <Text size="1" color="muted" weight="bold" transform="uppercase">{label}</Text>
        <Text size="8" weight="bold" color="cyan">{value}</Text>
      </Stack>
    </Card>
  );
}
