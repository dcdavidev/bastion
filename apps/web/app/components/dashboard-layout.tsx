import type { ReactNode } from 'react';

import { Outlet, useLocation, useNavigate } from 'react-router';

import {
  IconLayoutDashboard,
  IconLock,
  IconLogout,
  IconShieldExclamation,
  IconUsers,
  IconUserShield,
} from '@tabler/icons-react';

import {
  Avatar,
  Box,
  Button,
  Chip,
  Container,
  Flex,
  Stack,
  Text,
} from '@pittorica/react';

import { useAuth } from '../contexts/auth-context';

export default function DashboardLayout() {
  const { logout, isAdmin } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const version = import.meta.env.VITE_BASTION_VERSION || 'dev';

  return (
    <Box
      style={{
        display: 'grid',
        gridTemplateRows: '64px 1fr',
        gridTemplateColumns: '280px 1fr',
        height: '100vh',
        width: '100vw',
        backgroundColor: 'var(--pittorica-source-color)',
        overflow: 'hidden',
      }}
    >
      {/* APP BAR (Header) */}
      <Flex
        style={{
          gridColumn: '1 / span 2',
          borderBottom: '1px solid var(--pittorica-source-color)',
          backgroundColor: 'var(--pittorica-source-color)',
          zIndex: 10,
        }}
        px="4"
        align="center"
        justify="between"
      >
        <Flex align="center" gap="3" style={{ width: '264px' }}>
          <Avatar src="/static/logo/square.png" fallback="B" size="3" />
          <Text
            weight="bold"
            color="white"
            size="4"
            style={{ letterSpacing: '0.5px' }}
          >
            BASTION
          </Text>
        </Flex>
      </Flex>

      {/* SIDE NAV (Sidebar) */}
      <Flex
        p="4"
        direction="column"
        style={{
          backgroundColor: 'var(--pittorica-color-surface)',
          borderRight: '1px solid transparent',
        }}
      >
        <Stack gap="1" style={{ flex: 1 }}>
          <NavItem
            to="/dashboard"
            icon={<IconLayoutDashboard size={20} />}
            label="Overview"
            end
          />
          <NavItem
            to="/dashboard/clients"
            icon={<IconUsers size={20} />}
            label="Clients"
          />

          {isAdmin && (
            <>
              <Box mt="4" mb="1" px="4">
                <Text
                  size="1"
                  weight="bold"
                  color="muted"
                  style={{ textTransform: 'uppercase' }}
                >
                  Admin
                </Text>
              </Box>
              <NavItem
                to="/dashboard/collaborators"
                icon={<IconUserShield size={20} />}
                label="Collaborators"
              />
              <NavItem
                to="/dashboard/audit"
                icon={<IconShieldExclamation size={20} />}
                label="Audit Logs"
              />
            </>
          )}
        </Stack>

        <Box
          pt="4"
          style={{ borderTop: '1px solid var(--pittorica-color-border)' }}
        >
          <Stack gap="2" px="2">
            <Flex align="center" gap="3" p="2">
              <IconLock size={16} />
              <Text size="1" color="muted">
                E2EE Active
              </Text>
            </Flex>
            <Button
              variant="filled"
              onClick={handleLogout}
              color="red"
              size="sm"
              style={{ width: '100%', justifyContent: 'flex-start' }}
            >
              <IconLogout size={18} />
              <Text>Sign out</Text>
            </Button>
            <Flex justify="center" mt="2">
              <Chip variant="soft" color="source">
                v{version}
              </Chip>
            </Flex>
          </Stack>
        </Box>
      </Flex>

      {/* MAIN CONTENT */}
      <Box
        as="main"
        style={{
          overflow: 'auto',
          backgroundColor: 'var(--pittorica-color-surface)',
          paddingTop: '8px',
        }}
      >
        <Box
          style={{
            backgroundColor: 'var(--pittorica-color-surface-container-low)',
            minHeight: 'calc(100vh - 72px)',
            marginRight: '16px',
            marginBottom: '16px',
            borderRadius: '24px',
            border: '1px solid var(--pittorica-color-border)',
          }}
        >
          <Container>
            <Box p="8">
              <Outlet />
            </Box>
          </Container>
        </Box>
      </Box>
    </Box>
  );
}

function NavItem({
  to,
  icon,
  label,
  end = false,
}: {
  to: string;
  icon: ReactNode;
  label: string;
  end?: boolean;
}) {
  const navigate = useNavigate();
  const location = useLocation();

  const isActive = end
    ? location.pathname === to
    : location.pathname.startsWith(to);

  return (
    <Button
      variant={isActive ? 'tonal' : 'text'}
      onClick={() => navigate(to)}
      style={{
        justifyContent: 'flex-start',
        borderRadius: '0 100px 100px 0',
        paddingLeft: '24px',
        marginRight: '8px',
        height: '48px',
        backgroundColor: isActive
          ? 'rgba(var(--pittorica-color-source-rgb), 0.15)'
          : 'transparent',
        color: 'var(--pittorica-color-source)',
      }}
    >
      <Flex gap="4" align="center">
        <Box
          style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          {icon}
        </Box>
        <Text weight={isActive ? 'bold' : 'medium'} size="2">
          {label}
        </Text>
      </Flex>
    </Button>
  );
}
