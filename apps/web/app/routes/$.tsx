import { useNavigate } from 'react-router';

import { IconChevronLeft, IconHome, IconSearchOff } from '@tabler/icons-react';

import {
  Box,
  Button,
  Card,
  Container,
  Flex,
  Heading,
  Stack,
  Text,
} from '@pittorica/react';

export default function NotFound() {
  const navigate = useNavigate();

  return (
    <Box
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        background: 'var(--pittorica-color-surface)',
      }}
    >
      <Container maxWidth="sm">
        <Card
          p="8"
          style={{
            textAlign: 'center',
            boxShadow: 'var(--pittorica-shadow-lg)',
          }}
        >
          <Stack gap="6" align="center">
            <Box
              p="4"
              style={{
                backgroundColor: 'rgba(var(--pittorica-color-source-rgb), 0.1)',
                borderRadius: 'var(--pittorica-radius-full)',
                color: 'var(--pittorica-color-source)',
              }}
            >
              <IconSearchOff size={48} />
            </Box>

            <Stack gap="2">
              <Text
                size="9"
                weight="bold"
                style={{
                  lineHeight: 1,
                  opacity: 0.2,
                  fontFamily: 'var(--pittorica-font-code)',
                }}
              >
                404
              </Text>
              <Heading size="6" weight="bold">
                Access Denied: Page Missing
              </Heading>
              <Text color="muted" size="3">
                The vault section you are trying to access does not exist or has
                been unauthorized.
              </Text>
            </Stack>

            <Flex gap="3" style={{ width: '100%' }}>
              <Button
                variant="tonal"
                style={{ flex: 1 }}
                onClick={() => navigate(-1)}
              >
                <Flex align="center" gap="2">
                  <IconChevronLeft size={18} />
                  Go Back
                </Flex>
              </Button>
              <Button
                variant="filled"
                style={{ flex: 1 }}
                onClick={() => navigate('/')}
              >
                <Flex align="center" gap="2">
                  <IconHome size={18} />
                  Dashboard
                </Flex>
              </Button>
            </Flex>
          </Stack>
        </Card>
      </Container>
    </Box>
  );
}
