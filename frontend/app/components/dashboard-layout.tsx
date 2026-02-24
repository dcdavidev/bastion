import { useNavigate, Link, Outlet } from "react-router";
import { Box, Button, Text, Stack, Card } from "@pittorica/react";
import { useAuth } from "../contexts/auth-context";
import { LayoutDashboard, Users, ShieldAlert, LogOut } from "lucide-react";

export default function DashboardLayout() {
  const { logout, isAdmin } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <Box display="flex" minHeight="100vh" background="surface">
      {/* Sidebar */}
      <Box 
        width="260px" 
        padding="6" 
        borderRight="1px solid border" 
        background="surface-container"
        display="flex"
        flexDirection="column"
      >
        <Stack gap="8" flex="1">
          <Box>
            <Text size="5" weight="bold" color="cyan">BASTION</Text>
            <Text size="1" color="muted">Management Console</Text>
          </Box>

          <Stack gap="2">
            <NavItem to="/dashboard" icon={<LayoutDashboard size={18} />} label="Overview" />
            <NavItem to="/dashboard/clients" icon={<Users size={18} />} label="Clients" />
            {isAdmin && (
              <NavItem to="/dashboard/audit" icon={<ShieldAlert size={18} />} label="Audit Logs" />
            )}
          </Stack>
        </Stack>

        <Button 
          variant="ghost" 
          onClick={handleLogout}
          width="100%"
          justifyContent="flex-start"
        >
          <Stack direction="row" gap="3" alignItems="center">
            <LogOut size={18} />
            <Text>Logout</Text>
          </Stack>
        </Button>
      </Box>

      {/* Main Content */}
      <Box flex="1" padding="8" overflow="auto">
        <Outlet />
      </Box>
    </Box>
  );
}

function NavItem({ to, icon, label }: { to: string; icon: React.ReactNode; label: string }) {
  return (
    <Link to={to} style={{ textDecoration: 'none', color: 'inherit' }}>
      <Box 
        padding="3" 
        borderRadius="md" 
        className="hover:bg-accent/10 transition-colors"
        display="flex"
        alignItems="center"
      >
        <Stack direction="row" gap="3" alignItems="center">
          <Box color="cyan">{icon}</Box>
          <Text weight="medium">{label}</Text>
        </Stack>
      </Box>
    </Link>
  );
}
