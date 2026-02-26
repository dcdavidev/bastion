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
} from 'react-router';

import { PittoricaTheme, ToastProvider } from '@pittorica/react';

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
  let message = 'Oops!';
  let details = 'An unexpected error occurred.';
  let stack: string | undefined;

  if (isRouteErrorResponse(error)) {
    message = error.status === 404 ? '404' : 'Error';
    details =
      error.status === 404
        ? 'The requested page could not be found.'
        : error.statusText || details;
  } else if (import.meta.env.DEV && error && error instanceof Error) {
    details = error.message;
    stack = error.stack;
  }

  return (
    <main className="pt-16 p-4 container mx-auto">
      <h1>{message}</h1>
      <p>{details}</p>
      {stack && (
        <pre className="w-full p-4 overflow-x-auto">
          <code>{stack}</code>
        </pre>
      )}
    </main>
  );
}
