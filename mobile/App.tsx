import React from 'react';
import { Platform } from 'react-native';
import { NavigationContainer } from '@react-navigation/native';
import DrawerNavigator from './src/navigation/DrawerNavigator';
import { GestureHandlerRootView } from 'react-native-gesture-handler';
import { Client, Provider, cacheExchange, fetchExchange } from 'urql';
import 'react-native-url-polyfill/auto';

import { AuthProvider, useAuth } from './src/context/AuthContext';
import LoginScreen from './src/screens/LoginScreen';

const host = 'localhost';

function AppContent() {
  const { isLoggedIn, token } = useAuth();

  const client = React.useMemo(() => {
    return new Client({
      url: `http://${host}:3030/api/graphql`,
      exchanges: [cacheExchange, fetchExchange],
      preferGetMethod: false,
      fetchOptions: {
        headers: {
          Authorization: token ? `Bearer ${token}` : '',
        },
      },
      fetch: (url, options) => {
        return fetch(url, options);
      },
    });
  }, [token]);

  if (!isLoggedIn) {
    return <LoginScreen />;
  }

  return (
    <Provider value={client}>
      <NavigationContainer>
        <DrawerNavigator />
      </NavigationContainer>
    </Provider>
  );
}

import { ErrorBoundary } from './src/components/ErrorBoundary';

function App(): React.JSX.Element {
  return (
    <GestureHandlerRootView style={{ flex: 1 }}>
      <ErrorBoundary>
        <AuthProvider>
          <AppContent />
        </AuthProvider>
      </ErrorBoundary>
    </GestureHandlerRootView>
  );
}

export default App;
