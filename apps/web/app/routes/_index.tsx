import { useNavigate } from 'react-router';

import {
  IconArrowRight,
  IconLock,
  IconServerOff,
  IconUsers,
} from '@tabler/icons-react';

import {
  Avatar,
  Box,
  Button,
  Card,
  Container,
  Flex,
  Grid,
  Stack,
  Text,
} from '@pittorica/react';

export default function Home() {
  const navigate = useNavigate();

  return (
    <Box
      style={{
        minHeight: '100vh',
        background:
          'radial-gradient(circle at top right, #00bcd40d, transparent), radial-gradient(circle at bottom left, #00bcd40d, transparent)',
      }}
    >
      <Container maxWidth="lg" p="8">
        <Stack
          gap="4"
          align="center"
          style={{ textAlign: 'center', paddingTop: '10vh' }}
        >
          <Stack gap="4" align="center">
            <Avatar src="/static/logo/square.png" fallback="B" size="9" />
            <Stack gap="1">
              <Text size="5" color="source" weight="medium">
                The Secure Fortress for Your Secrets
              </Text>
            </Stack>
          </Stack>

          <Text
            size="4"
            color="muted"
            style={{ maxWidth: '600px', lineHeight: '1.6' }}
          >
            A self-hosted, end-to-end encrypted vault designed for developers
            and teams. Manage multiple client secrets with zero-knowledge
            architecture and a blind-backend approach.
          </Text>

          <Flex gap="4">
            <Button
              size="lg"
              variant="filled"
              onClick={() => navigate('/login')}
              style={{ minWidth: '200px' }}
            >
              <Flex gap="2" align="center">
                Unlock Vault
                <IconArrowRight size={18} />
              </Flex>
            </Button>
            <Button
              size="lg"
              variant="tonal"
              as="a"
              href="https://github.com/dcdavidev/bastion"
              target="_blank"
            >
              View Source
            </Button>
          </Flex>

          <Grid columns="3" gap="6" style={{ width: '100%', marginTop: '4vh' }}>
            <FeatureCard
              icon={<IconLock size={24} />}
              title="E2E Encrypted"
              description="Data is encrypted in your browser. We never see your master key or plaintext secrets."
            />
            <FeatureCard
              icon={<IconServerOff size={24} />}
              title="Blind Backend"
              description="Our Go backend only stores encrypted blobs. Even with full db access, your data remains safe."
            />
            <FeatureCard
              icon={<IconUsers size={24} />}
              title="Multi-Client"
              description="Organize secrets by clients and projects with isolated encryption keys for each."
            />
          </Grid>
        </Stack>
      </Container>

      <Box
        p="8"
        style={{
          borderTop: '1px solid var(--pittorica-color-border)',
          marginTop: '10vh',
          textAlign: 'center',
        }}
      >
        <Text size="1" color="muted">
          &copy; {new Date().getFullYear()} Bastion E2EE. Built with Go, React,
          and Pittorica.
        </Text>
      </Box>
    </Box>
  );
}

function FeatureCard({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
}) {
  return (
    <Card p="6">
      <Stack gap="4" align="center">
        <Box color="cyan">{icon}</Box>
        <Stack gap="2">
          <Text weight="bold" size="3">
            {title}
          </Text>
          <Text size="2" color="muted" style={{ lineHeight: '1.4' }}>
            {description}
          </Text>
        </Stack>
      </Stack>
    </Card>
  );
}
