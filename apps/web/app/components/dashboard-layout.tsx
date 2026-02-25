import type { ReactNode } from 'react';

import { Link, Outlet, useNavigate } from 'react-router';

import {
  IconLayoutDashboard,
  IconLogout,
  IconShieldExclamation,
  IconUsers,
} from '@tabler/icons-react';

import { 
  Flex, 
  Container, 
  Button, 
  Stack, 
  Text, 
  Box, 
  IconButton,
  Avatar
} from '@pittorica/react';

import { useAuth } from '../contexts/auth-context';

export default function DashboardLayout() {
  const { logout, isAdmin } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <Flex
      height="100vh"
      style={{ backgroundColor: 'var(--pittorica-color-surface)' }}
    >
      {/* Sidebar */}
      <Flex
        width="260px"
        p="6"
        direction="column"
        style={{
          borderRight: '1px solid var(--pittorica-color-border)',
          backgroundColor: 'var(--pittorica-color-surface-container)',
        }}
      >
        <Stack gap="8" style={{ flex: 1 }}>
          <Flex gap="3" align="center">
            <Avatar 
              src="/static/logo/square.png" 
              fallback="B" 
              size="md" 
            />
            <Box>
              <Text size="4" weight="bold" color="cyan">
                BASTION
              </Text>
              <Text size="1" color="muted">
                v1.0.0-alpha
              </Text>
            </Box>
          </Flex>

          <Stack gap="2">
            <NavItem
              to="/dashboard"
              icon={<IconLayoutDashboard size={18} />}
              label="Overview"
            />
            <NavItem
              to="/dashboard/clients"
              icon={<IconUsers size={18} />}
              label="Clients"
            />
            {isAdmin && (
              <NavItem
                to="/dashboard/audit"
                icon={<IconShieldExclamation size={18} />}
                label="Audit Logs"
              />
            )}
          </Stack>
        </Stack>

        <Button
          variant="text"
          onClick={handleLogout}
          width="100%"
          style={{ justifyContent: 'flex-start' }}
        >
          <Stack direction="row" gap="3" align="center">
            <IconLogout size={18} />
            <Text>Logout</Text>
          </Stack>
        </Button>
      </Flex>

      {/* Main Content Area */}
      <Container style={{ flex: 1, overflow: 'auto' }}>
        <Flex direction="column" p="8">
          <Outlet />
        </Flex>
      </Container>
    </Flex>
  );
}

function NavItem({
  to,
  icon,
  label,
}: {
  to: string;
  icon: ReactNode;
  label: string;
}) {
  return (
    <Link to={to} style={{ textDecoration: 'none', color: 'inherit' }}>
      <Flex
        p="3"
        align="center"
        className="hover:bg-accent/10 transition-colors"
        style={{ borderRadius: 'var(--pittorica-radius-md)' }}
      >
        <Stack direction="row" gap="3" align="center">
          <IconButton variant="text" color="cyan">
            {icon}
          </IconButton>
          <Text weight="medium">{label}</Text>
        </Stack>
      </Flex>
    </Link>
  );
}
