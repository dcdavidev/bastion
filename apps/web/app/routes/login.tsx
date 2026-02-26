import type { ChangeEvent, SyntheticEvent } from 'react';
import { useState } from 'react';

import { useNavigate } from 'react-router';

import { IconEye, IconEyeOff } from '@tabler/icons-react';

import {
  Avatar,
  Button,
  Card,
  Flex,
  IconButton,
  Stack,
  Text,
  TextField,
  toast,
} from '@pittorica/react';

import { useAuth } from '../contexts/auth-context';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const navigate = useNavigate();
  const { setToken } = useAuth();

  // Calculate form validity (password is always required, email is optional for admin fallback)
  const isFormValid = password.trim().length > 0;

  async function handleSubmit(e: SyntheticEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!isFormValid) return;

    setLoading(true);
    setError('');

    try {
      const response = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: email.trim() || undefined,
          password,
        }),
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
            <Avatar src="/static/logo/square.png" fallback="B" size="9" />
            <Stack gap="0" align="center">
              <Text size="7" weight="bold">
                BASTION
              </Text>
              <Text size="2" color="muted">
                Secure E2EE Secrets Vault
              </Text>
            </Stack>
          </Flex>

          <form onSubmit={handleSubmit}>
            <Stack gap="4">
              <TextField.Root label="Email (optional)" size="md">
                <TextField.Input
                  placeholder="Leave empty for Admin fallback"
                  value={email}
                  onChange={(e: ChangeEvent<HTMLInputElement>) =>
                    setEmail(e.target.value)
                  }
                />
              </TextField.Root>

              <TextField.Root label="Password" size="md">
                <TextField.Input
                  type={showPassword ? 'text' : 'password'}
                  required
                  placeholder="••••••••"
                  value={password}
                  onChange={(e: ChangeEvent<HTMLInputElement>) =>
                    setPassword(e.target.value)
                  }
                />
                <TextField.Slot
                  style={{ cursor: 'pointer', paddingRight: '8px' }}
                >
                  <IconButton
                    type="button"
                    variant="text"
                    size="2"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? (
                      <IconEyeOff size={20} />
                    ) : (
                      <IconEye size={20} />
                    )}
                  </IconButton>
                </TextField.Slot>
              </TextField.Root>

              {error && (
                <Text size="2" color="red">
                  {error}
                </Text>
              )}

              <Button
                type="submit"
                variant="filled"
                disabled={loading || !isFormValid}
                width="100%"
                size="md"
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
