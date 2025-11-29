import React, { useEffect, useState } from 'react';
import {
  StatusBar,
  StyleSheet,
  Text,
  useColorScheme,
  View,
  Platform,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

function App(): React.JSX.Element {
  const isDarkMode = useColorScheme() === 'dark';
  const [status, setStatus] = useState<string>('Checking...');
  const [backendUrl, setBackendUrl] = useState<string>('');

  useEffect(() => {
    const host = Platform.OS === 'android' ? '10.0.2.2' : 'localhost';
    const url = `http://${host}:3030/health`;
    setBackendUrl(url);

    fetch(url)
      .then(res => {
        if (res.ok) return 'Backend Connected! Hurrey!!';
        return `Error: ${res.status}`;
      })
      .catch(err => `Fetch Error: ${err.message}`)
      .then(setStatus);
  }, []);

  return (
    <SafeAreaView style={styles.container}>
      <StatusBar barStyle={isDarkMode ? 'light-content' : 'dark-content'} />
      <View style={styles.content}>
        <Text style={styles.text}>GoAlert Mobile</Text>
        <Text style={styles.text}>Backend URL: {backendUrl}</Text>
        <Text style={styles.status}>{status}</Text>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#ffffffff',
  },
  content: {
    padding: 20,
    alignItems: 'center',
  },
  text: {
    fontSize: 18,
    marginBottom: 10,
    color: '#333',
  },
  status: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#007AFF',
  },
});

export default App;
