import React, { useEffect, useState } from 'react';
import { View, Text, TextInput, Button, StyleSheet, ActivityIndicator, Alert } from 'react-native';
import { useAuth } from '../context/AuthContext';

const host = 'localhost'; // Same as App.tsx
const BASE_URL = `http://${host}:3030`;

export default function LoginScreen() {
  const { login } = useAuth();
  const [providers, setProviders] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  useEffect(() => {
    fetch(`${BASE_URL}/api/v2/identity/providers`)
      .then((res) => res.json())
      .then((data) => {
        setProviders(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error('Failed to fetch providers:', err);
        setLoading(false);
      });
  }, []);

  const handleLogin = async (providerId: string, url: string) => {
    const formData = new URLSearchParams();
    formData.append('username', username);
    formData.append('password', password);
    // Add other fields if necessary based on provider.Fields

    const fullUrl = BASE_URL + url + '?noRedirect=1';

    try {
      const res = await fetch(fullUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
          'Referer': `http://${host}:3030`,
        },
        body: formData.toString(),
      });

      const text = await res.text();

      if (res.ok) {
        login(text);
      } else {
        Alert.alert('Login Failed', `Status: ${res.status}\n${text}`);
      }
    } catch (err) {
      console.error('Login error:', err);
      Alert.alert('Login Error', 'An error occurred during login');
    }
  };

  if (loading) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" />
      </View>
    );
  }

  // Assuming basic auth provider for now, or the first one found
  const basicProvider = providers.find(p => p.ID === 'basic');

  if (!basicProvider) {
    return (
      <View style={styles.container}>
        <Text>No Basic Auth Provider found</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <Text style={styles.title}>GoAlert Login</Text>
      <TextInput
        style={styles.input}
        placeholder="Username"
        value={username}
        onChangeText={setUsername}
        autoCapitalize="none"
      />
      <TextInput
        style={styles.input}
        placeholder="Password"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
      />
      <Button title="Login" onPress={() => handleLogin(basicProvider.ID, basicProvider.URL)} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    padding: 20,
    backgroundColor: '#fff',
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 20,
    textAlign: 'center',
  },
  input: {
    borderWidth: 1,
    borderColor: '#ccc',
    padding: 10,
    marginBottom: 10,
    borderRadius: 5,
  },
});
