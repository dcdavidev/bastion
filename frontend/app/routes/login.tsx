import type { ChangeEvent, SyntheticEvent } from 'react';
import { useState } from 'react';

import { useNavigate } from 'react-router';

import {
  Avatar,
  Button,
  Card,
  Flex,
  Stack,
  Text,
  TextField,
  toast,
} from '@pittorica/react';

import { useAuth } from '../contexts/auth-context';

export default function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const { setToken } = useAuth();

  async function handleSubmit(e: SyntheticEvent<HTMLFormElement>) {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) {
        throw new Error('Invalid credentials');
      }

      const data = await response.json();
      setToken(data.token);
      toast({
        title: 'Access granted',
        description: 'Welcome back!',
        color: 'teal',
      });
      navigate('/dashboard');
    } catch (error: unknown) {
      const message =
        error instanceof Error ? error.message : 'Something went wrong';
      setError(message);
      toast({
        title: 'Authentication failed',
        description: message,
        color: 'red',
      });
    } finally {
      setLoading(false);
    }
  }

  return (
    <Flex align="center" justify="center" height="100vh" p="4">
      <Card p="6" style={{ width: '100%', maxWidth: '400px' }}>
        <Stack gap="6">
          <Flex direction="column" gap="3" justify="center" align="center">
            <Avatar 
              src="/static/logo/square.png" 
              fallback="B" 
              size="lg" 
            />
            <Stack gap="0" align="center">
              <Text size="7" weight="bold" color="cyan">
                BASTION
              </Text>
              <Text size="2" color="muted">
                Secure E2EE Secrets Vault
              </Text>
            </Stack>
          </Flex>

          <form onSubmit={handleSubmit}>
            <Stack gap="4">
              <TextField.Root label="Username (optional)">
                <TextField.Input
                  placeholder="Leave empty for Admin"
                  value={username}
                  onChange={(e: ChangeEvent<HTMLInputElement>) =>
                    setUsername(e.target.value)
                  }
                />
              </TextField.Root>

              <TextField.Root label="Password">
                <TextField.Input
                  type="password"
                  required
                  placeholder="••••••••"
                  value={password}
                  onChange={(e: ChangeEvent<HTMLInputElement>) =>
                    setPassword(e.target.value)
                  }
                />
              </TextField.Root>

              {error && (
                <Text size="2" color="red">
                  {error}
                </Text>
              )}

              <Button
                type="submit"
                variant="filled"
                disabled={loading}
                width="100%"
              >
                {loading ? 'Unlocking...' : 'Unlock Vault'}
              </Button>
            </Stack>
          </form>
        </Stack>
      </Card>
    </Flex>
  );
}
