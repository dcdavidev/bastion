import 'pittorica';
import '@fontsource/inknut-antiqua';
import '@fontsource-variable/inter';
import '@fontsource/momo-trust-display';
import '@fontsource-variable/kode-mono';

import type { ReactNode } from 'react';

import {
  isRouteErrorResponse,
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  useNavigate,
} from 'react-router';

import {
  IconAlertTriangle,
  IconChevronLeft,
  IconHome,
  IconRefresh,
  IconShieldOff,
} from '@tabler/icons-react';

import {
  Box,
  Button,
  Card,
  Code,
  Container,
  Flex,
  Heading,
  PittoricaTheme,
  Stack,
  Text,
  ToastProvider,
} from '@pittorica/react';

import './app.css';

import type { Route } from './+types/root';
import { AuthProvider } from './contexts/auth-context';

export const links: Route.LinksFunction = () => [
  {
    rel: 'icon',
    type: 'image/png',
    href: '/favicon-96x96.png',
    sizes: '96x96',
  },
  { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' },
  { rel: 'shortcut icon', href: '/favicon.ico' },
  { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' },
  { rel: 'manifest', href: '/site.webmanifest' },
];

export const meta: Route.MetaFunction = () => [
  { title: 'Bastion' },
  {
    name: 'description',
    content:
      'Bastion is a single-user, open-source E2EE secrets vault built with Go. It provides a secure, self-hosted fortress to manage multiple client secrets via a powerful CLI and dashboard, ensuring data stays private with blind-backend architecture.',
  },
  { name: 'apple-mobile-web-app-title', content: 'Bastion' },
  { property: 'og:title', content: 'Bastion' },
  { property: 'og:description', content: 'Secure E2EE secrets vault' },
  { property: 'og:image', content: '/og-image.png' },
  { name: 'twitter:card', content: 'summary_large_image' },
  { name: 'twitter:image', content: '/og-image.png' },
];

export function Layout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta
          name="viewport"
          content="width=device-width, initial-scale=1, user-scalable=no"
        />
        <Meta />
        <Links />
      </head>
      <body className="pittorica-theme" data-appearance="light">
        <PittoricaTheme sourceColor="#388697" appearance="light">
          <AuthProvider>{children}</AuthProvider>
          <ToastProvider />
          <ScrollRestoration />
        </PittoricaTheme>
        <Scripts />
      </body>
    </html>
  );
}

export default function App() {
  return <Outlet />;
}

export function ErrorBoundary({ error }: Route.ErrorBoundaryProps) {
  const navigate = useNavigate();
  let status = '500';
  let title = 'Internal Server Error';
  let message =
    'The Bastion has encountered an unexpected security breach in our code.';
  let stack: string | undefined;

  if (isRouteErrorResponse(error)) {
    status = error.status.toString();
    if (error.status === 404) {
      title = 'Page Not Found';
      message =
        "The vault section you're looking for doesn't exist or has been moved.";
    } else {
      title = error.statusText || title;
    }
  } else if (error instanceof Error) {
    message = error.message;
    stack = error.stack;
  }

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
              {status === '404' ? (
                <IconShieldOff size={48} />
              ) : (
                <IconAlertTriangle size={48} />
              )}
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
                {status}
              </Text>
              <Heading size="6" weight="bold">
                {title}
              </Heading>
              <Text color="muted" size="3">
                {message}
              </Text>
            </Stack>

            {stack && (
              <Box style={{ width: '100%', textAlign: 'left' }}>
                <Text
                  size="1"
                  weight="bold"
                  color="muted"
                  mb="2"
                  style={{ textTransform: 'uppercase' }}
                >
                  Stack Trace
                </Text>
                <Box
                  p="4"
                  style={{
                    backgroundColor: 'var(--pittorica-color-surface-container)',
                    borderRadius: 'var(--pittorica-radius-md)',
                    maxHeight: '200px',
                    overflow: 'auto',
                  }}
                >
                  <Code
                    size="1"
                    style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-all' }}
                  >
                    {stack}
                  </Code>
                </Box>
              </Box>
            )}

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
                onClick={() => (globalThis.location.href = '/')}
              >
                <Flex align="center" gap="2">
                  <IconHome size={18} />
                  Home
                </Flex>
              </Button>
            </Flex>

            <Button
              variant="text"
              size="sm"
              onClick={() => globalThis.location.reload()}
            >
              <Flex align="center" gap="2">
                <IconRefresh size={14} />
                Try Reloading
              </Flex>
            </Button>
          </Stack>
        </Card>
      </Container>
    </Box>
  );
}
