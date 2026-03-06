import { useCallback } from 'react';

import {
  Outlet,
  redirect,
  useLoaderData,
  useLocation,
  useNavigate,
} from 'react-router';

import axios from 'axios';

import Cookies from 'js-cookie';
import {
  Box,
  Button,
  Card,
  Divider,
  Flex,
  Text,
  toast,
} from '@pittorica/react';
import { startRegistration } from '@simplewebauthn/browser';

import { api } from '../configs/api';

/**
 * clientLoader to protect all vault routes.
 * Verifies the JWT token from localStorage or cookies.
 */
export async function clientLoader() {
  // Try localStorage first (more reliable in SPA), then cookie
  let token = null;
  if (globalThis.window !== undefined) {
    token = localStorage.getItem('bastion_token');
  }

  if (!token) {
    token = Cookies.get('bastion_session');
  }

  if (!token) {
    console.warn('No authentication token found, redirecting to login');
    return redirect('/login?next=/vault');
  }

  try {
    const res = await api.get('/auth/me', {
      headers: { Authorization: `Bearer ${token}` },
    });
    return { authenticated: true, user: res.data };
  } catch (error: unknown) {
    console.error('Vault Auth Error:', error);

    if (axios.isAxiosError(error) && error.response?.status === 401) {
      if (globalThis.window !== undefined) {
        localStorage.removeItem('bastion_token');
      }
      Cookies.remove('bastion_session', { path: '/' });
    }

    return redirect('/login?next=/vault');
  }
}

export function shouldRevalidate() {
  return false;
}

interface AuthenticatedUser {
  id: string;
  username: string;
  email: string;
  role: string;
}

interface VaultLoaderData {
  authenticated: boolean;
  user: AuthenticatedUser;
}

export default function VaultLayout() {
  const loaderData = useLoaderData<VaultLoaderData>();
  const navigate = useNavigate();

  if (!loaderData || !loaderData.user) {
    return null;
  }

  const { user } = loaderData;

  const handleRegisterPasskey = useCallback(async () => {
    try {
      const optionsRes = await api.get('/auth/passkey/register/begin');
      const options = optionsRes.data;

      const regRes = await startRegistration(options);

      await api.post('/auth/passkey/register/finish', regRes);

      toast({
        title: 'Passkey registered',
        description: `Device registered for ${user.username}. You can now use this Passkey for future logins.`,
        color: 'teal',
      });
    } catch (error: unknown) {
      console.error('Passkey Registration Error:', error);
      toast({
        title: 'Error',
        description: 'Failed to register Passkey.',
        color: 'red',
      });
    }
  }, [user.username]);

  return (
    <Flex
      direction="row"
      style={{ minHeight: 'calc(100vh - 64px)', marginTop: '64px' }}
    >
      {/* Sidebar - Fixed on the left */}
      <Box
        as="aside"
        style={{
          width: '260px',
          height: 'calc(100vh - 64px)',
          position: 'fixed',
          left: 0,
          backgroundColor: 'var(--pittorica-surface-2)',
          borderRight: '1px solid var(--pittorica-outline-variant)',
          overflowY: 'auto',
          zIndex: 5,
        }}
      >
        <Flex direction={'column'} p="4">
          <Text
            weight="bold"
            size="2"
            color="muted"
            mb="4"
            style={{ textTransform: 'uppercase', letterSpacing: '1px' }}
          >
            Control Panel
          </Text>

          <Flex direction="column" gap="2">
            <VaultLink to="/vault" label="Overview" />
            <VaultLink to="/vault/clients" label="Clients" />
            <VaultLink to="/vault/audit" label="Audit Logs" />
            <VaultLink to="/vault/collaborators" label="Collaborators" />

            <Divider my="4" />

            <Button
              variant="tonal"
              onClick={handleRegisterPasskey}
              style={{ width: '100%', justifyContent: 'start' }}
            >
              Register Passkey
            </Button>

            <Button
              variant="tonal"
              color="red"
              onClick={() => {
                if (globalThis.window !== undefined) {
                  localStorage.removeItem('bastion_token');
                }
                Cookies.remove('bastion_session', { path: '/' });
                navigate('/login');
              }}
              style={{ width: '100%', justifyContent: 'start' }}
            >
              Logout
            </Button>
          </Flex>
        </Flex>
      </Box>

      {/* Main Content - Pushed to the right by the sidebar width */}
      <Box
        as="main"
        style={{
          flex: 1,
          marginLeft: '260px',
          padding: '24px',
          backgroundColor: 'var(--pittorica-surface-1)',
        }}
      >
        <Card p="6" style={{ minHeight: '100%', borderRadius: '12px' }}>
          <Outlet />
        </Card>
      </Box>
    </Flex>
  );
}

function VaultLink({ to, label }: { to: string; label: string }) {
  const navigate = useNavigate();
  const location = useLocation();
  const isActive = location.pathname === to;

  return (
    <Button
      variant={isActive ? 'filled' : 'tonal'}
      onClick={() => navigate(to)}
      style={{
        width: '100%',
        justifyContent: 'start',
        textAlign: 'left',
      }}
    >
      {label}
    </Button>
  );
}
