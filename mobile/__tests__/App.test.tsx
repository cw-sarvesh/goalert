/**
 * @format
 */

import React from 'react';
import ReactTestRenderer from 'react-test-renderer';
import App from '../App';

// Mock Reanimated
jest.mock('react-native-reanimated', () => {
  const Reanimated = require('react-native-reanimated/mock');
  Reanimated.default.call = () => {};
  return Reanimated;
});

// Mock Gesture Handler
jest.mock('react-native-gesture-handler', () => {
  return {
    State: {
      UNDETERMINED: 0,
      BEGAN: 1,
      ACTIVE: 2,
      END: 3,
      FAILED: 4,
      CANCELLED: 5,
    },
    PanGestureHandler: 'PanGestureHandler',
    TapGestureHandler: 'TapGestureHandler',
    ScrollView: 'ScrollView',
    Switch: 'Switch',
    TextInput: 'TextInput',
    DrawerLayoutAndroid: 'DrawerLayoutAndroid',
    FlatList: 'FlatList',
  };
});

// Mock Screens
jest.mock('react-native-screens', () => {
  return {
    enableScreens: jest.fn(),
  };
});

// Mock Drawer Navigator
jest.mock('@react-navigation/drawer', () => {
  const React = require('react');
  return {
    createDrawerNavigator: () => {
      return {
        Navigator: ({ children }: { children: React.ReactNode }) => <>{children}</>,
        Screen: ({ children }: { children: React.ReactNode }) => <>{children}</>,
      };
    },
  };
});

test('renders correctly', async () => {
  let renderer;
  await ReactTestRenderer.act(async () => {
    renderer = ReactTestRenderer.create(<App />);
  });
  expect(renderer).toBeDefined();
});
