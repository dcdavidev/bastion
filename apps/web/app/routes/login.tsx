import type { ChangeEvent, SyntheticEvent } from 'react';
import { useState } from 'react';

import { useNavigate } from 'react-router';

import { IconEye, IconEyeOff, IconFingerprint } from '@tabler/icons-react';

import {
  Avatar,
  Button,
  Card,
  Divider,
  Flex,
  IconButton,
  Stack,
  Text,
  TextField,
  toast,
} from '@pittorica/react';
import { startAuthentication } from '@simplewebauthn/browser';
import Cookies from 'js-cookie';

import { useAuth } from '../contexts/auth-context';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [passkeyLoading, setPasskeyLoading] = useState(false);
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
      Cookies.set('bastion_session', data.token, { expires: 1, path: '/' });
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

  async function handlePasskeyLogin() {
    if (!email.trim()) {
      setError('Please enter your email to use Passkey');
      return;
    }

    setPasskeyLoading(true);
    setError('');

    try {
      // 1. Get options from server
      const beginResp = await fetch(
        `/api/v1/auth/passkey/login/begin?email=${encodeURIComponent(email)}`
      );
      if (!beginResp.ok) throw new Error('Failed to start Passkey login');

      const options = await beginResp.json();

      // 2. Browser authentication
      const assertion = await startAuthentication(options);

      // 3. Verify assertion on server
      const finishResp = await fetch(
        `/api/v1/auth/passkey/login/finish?email=${encodeURIComponent(email)}`,
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(assertion),
        }
      );

      if (!finishResp.ok) throw new Error('Passkey verification failed');

      const data = await finishResp.json();
      setToken(data.token);
      Cookies.set('bastion_session', data.token, { expires: 1, path: '/' });
      toast({
        title: 'Access granted',
        description: 'Signed in with Passkey',
        color: 'teal',
      });
      navigate('/dashboard');
    } catch (error: unknown) {
      console.error(error);
      const message =
        error instanceof Error ? error.message : 'Passkey login failed';
      setError(message);
    } finally {
      setPasskeyLoading(false);
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

          <Stack gap="4">
            <TextField.Root label="Email" size="md">
              <TextField.Input
                placeholder="Required for Passkey / Admin fallback"
                value={email}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setEmail(e.target.value)
                }
              />
            </TextField.Root>

            <form onSubmit={handleSubmit}>
              <Stack gap="4">
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
                  disabled={loading || passkeyLoading || !isFormValid}
                  width="100%"
                  size="md"
                >
                  {loading ? 'Unlocking...' : 'Unlock Vault'}
                </Button>
              </Stack>
            </form>

            <Flex align="center" gap="3">
              <Divider style={{ flex: 1 }} />
              <Text size="1" color="muted">
                OR
              </Text>
              <Divider style={{ flex: 1 }} />
            </Flex>

            <Button
              variant="tonal"
              width="100%"
              size="md"
              disabled={loading || passkeyLoading}
              onClick={handlePasskeyLogin}
            >
              <Flex align="center" gap="2">
                <IconFingerprint size={20} />
                <Text>
                  {passkeyLoading
                    ? 'Authenticating...'
                    : 'Sign in with Passkey'}
                </Text>
              </Flex>
            </Button>
          </Stack>
        </Stack>
      </Card>
    </Flex>
  );
}
