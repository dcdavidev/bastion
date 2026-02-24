import { useState } from "react";
import { useNavigate } from "react-router";
import { Box, Button, Input, Text, Stack, Card } from "@pittorica/react";

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const response = await fetch("/api/v1/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) {
        throw new Error("Invalid credentials");
      }

      const data = await response.json();
      localStorage.setItem("bastion_token", data.token);
      navigate("/dashboard");
    } catch (err: any) {
      setError(err.message || "Something went wrong");
    } finally {
      setLoading(false);
    }
  }

  return (
    <Box 
      display="flex" 
      alignItems="center" 
      justifyContent="center" 
      minHeight="100vh"
      padding="4"
    >
      <Card width="100%" maxWidth="400px" padding="6">
        <Stack gap="6">
          <Box textAlign="center">
            <Text size="7" weight="bold" color="cyan">BASTION</Text>
            <Text size="2" color="muted">Secure E2EE Secrets Vault</Text>
          </Box>

          <form onSubmit={handleSubmit}>
            <Stack gap="4">
              <Stack gap="2">
                <Text size="2" weight="medium">Username (optional)</Text>
                <Input 
                  placeholder="Leave empty for Admin" 
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                />
              </Stack>

              <Stack gap="2">
                <Text size="2" weight="medium">Password</Text>
                <Input 
                  type="password" 
                  required
                  placeholder="••••••••" 
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
              </Stack>

              {error && (
                <Text size="2" color="red">{error}</Text>
              )}

              <Button 
                type="submit" 
                variant="primary" 
                loading={loading}
                width="100%"
              >
                Unlock Vault
              </Button>
            </Stack>
          </form>
        </Stack>
      </Card>
    </Box>
  );
}
