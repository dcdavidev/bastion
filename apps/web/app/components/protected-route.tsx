import { Navigate, Outlet } from "react-router";
import { useAuth } from "../contexts/auth-context";

export default function ProtectedRoute() {
  const { token } = useAuth();

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
}
