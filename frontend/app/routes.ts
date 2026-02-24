import { 
  type RouteConfig, 
  index, 
  layout, 
  route 
} from "@react-router/dev/routes";

export default [
  // Public Routes
  route("login", "routes/login.tsx"),

  // Protected Routes
  layout("components/protected-route.tsx", [
    layout("components/dashboard-layout.tsx", [
      route("dashboard", "routes/dashboard/overview.tsx"),
      route("dashboard/clients", "routes/dashboard/clients.tsx"),
      route("dashboard/clients/:clientId", "routes/dashboard/projects.tsx"),
      route("dashboard/projects/:projectId", "routes/dashboard/secrets.tsx"),
      route("dashboard/audit", "routes/dashboard/audit.tsx"),
      route("dashboard/collaborators", "routes/dashboard/collaborators.tsx"),
    ]),
  ]),

  // Default redirect
  index("routes/_index.tsx"),
] satisfies RouteConfig;
