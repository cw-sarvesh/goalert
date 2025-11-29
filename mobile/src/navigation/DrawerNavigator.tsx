import React from 'react';
import { createDrawerNavigator } from '@react-navigation/drawer';
import AlertsScreen from '../screens/AlertsScreen';
import PlaceholderScreen from '../screens/PlaceholderScreen';
import { useTheme } from '../context/ThemeContext';

const Drawer = createDrawerNavigator();

export default function DrawerNavigator() {
  const { colors, isDark } = useTheme();

  return (
    <Drawer.Navigator
      screenOptions={{
        headerStyle: {
          backgroundColor: colors.primary,
        },
        headerTintColor: '#fff',
        drawerStyle: {
          backgroundColor: colors.card,
          width: '75%',
        },
        drawerActiveTintColor: colors.primary,
        drawerInactiveTintColor: colors.text,
        sceneContainerStyle: {
          backgroundColor: colors.background,
        },
      }}
    >
      <Drawer.Screen name="Alerts" component={AlertsScreen} />
      <Drawer.Screen name="Rotations" component={PlaceholderScreen} />
      <Drawer.Screen name="Schedules" component={PlaceholderScreen} />
      <Drawer.Screen name="Escalation Policies" component={PlaceholderScreen} />
      <Drawer.Screen name="Services" component={PlaceholderScreen} />
      <Drawer.Screen name="Users" component={PlaceholderScreen} />
      <Drawer.Screen name="Admin" component={PlaceholderScreen} />
      <Drawer.Screen name="Wizard" component={PlaceholderScreen} />
    </Drawer.Navigator>
  );
}
