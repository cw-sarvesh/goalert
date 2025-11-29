import React, { useState } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, ActivityIndicator } from 'react-native';
import { useQuery, gql } from 'urql';
import { useAuth } from '../context/AuthContext';
import { useTheme } from '../context/ThemeContext';

const USER_QUERY = gql`
  query CurrentUser {
    user {
      id
      name
      role
    }
  }
`;

export default function UserProfileMenu({ visible, onClose }: { visible: boolean; onClose: () => void }) {
  const { logout } = useAuth();
  const { theme, setTheme, colors } = useTheme();
  const [result] = useQuery({ query: USER_QUERY });
  const { data, fetching } = result;

  const userName = data?.user?.name || 'User';
  const firstName = userName.split(' ')[0];

  const handleLogout = () => {
    onClose();
    logout();
  };

  if (!visible) return null;

  return (
    <View style={styles.modalOverlay}>
      <TouchableOpacity style={styles.overlayBackground} activeOpacity={1} onPress={onClose} />
      <View style={[styles.menuContainer, { backgroundColor: colors.card }]}>
        <View style={styles.header}>
          {fetching ? (
            <ActivityIndicator size="small" color={colors.text} />
          ) : (
            <Text style={[styles.greeting, { color: colors.text }]}>Hello, {firstName}!</Text>
          )}
        </View>

        <TouchableOpacity style={styles.menuItem} onPress={() => console.log('Manage Profile')}>
          <Text style={[styles.menuText, { color: colors.text }]}>Manage Profile</Text>
        </TouchableOpacity>

        <View style={[styles.divider, { backgroundColor: colors.border }]} />

        <Text style={[styles.sectionHeader, { color: colors.subText }]}>Appearance</Text>
        <View style={styles.appearanceRow}>
          <TouchableOpacity
            style={[styles.appearanceBtn, { borderColor: colors.border }, theme === 'light' && styles.activeAppearance]}
            onPress={() => setTheme('light')}
          >
            <Text style={[styles.btnText, { color: colors.text }, theme === 'light' && styles.activeText]}>Light</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={[styles.appearanceBtn, { borderColor: colors.border }, theme === 'system' && styles.activeAppearance]}
            onPress={() => setTheme('system')}
          >
            <Text style={[styles.btnText, { color: colors.text }, theme === 'system' && styles.activeText]}>System</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={[styles.appearanceBtn, { borderColor: colors.border }, theme === 'dark' && styles.activeAppearance]}
            onPress={() => setTheme('dark')}
          >
            <Text style={[styles.btnText, { color: colors.text }, theme === 'dark' && styles.activeText]}>Dark</Text>
          </TouchableOpacity>
        </View>

        <View style={[styles.divider, { backgroundColor: colors.border }]} />

        <TouchableOpacity style={styles.menuItem} onPress={handleLogout}>
          <Text style={[styles.menuText, styles.logoutText]}>Logout</Text>
        </TouchableOpacity>

        <TouchableOpacity style={[styles.closeButton, { backgroundColor: colors.border }]} onPress={onClose}>
          <Text style={[styles.closeButtonText, { color: colors.text }]}>Close</Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  iconButton: {
    padding: 8,
  },
  iconText: {
    fontSize: 24,
    color: 'black',
  },
  modalOverlay: {
    position: 'absolute',
    top: 0,
    bottom: 0,
    left: 0,
    right: 0,
    justifyContent: 'flex-end',
    zIndex: 1000,
  },
  overlayBackground: {
    position: 'absolute',
    top: 0,
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: 'rgba(0,0,0,0.5)',
  },
  menuContainer: {
    backgroundColor: 'white',
    borderTopLeftRadius: 20,
    borderTopRightRadius: 20,
    padding: 20,
    width: '100%',
    elevation: 5,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: -2 },
    shadowOpacity: 0.25,
    shadowRadius: 3.84,
  },
  header: {
    marginBottom: 20,
    alignItems: 'center',
  },
  greeting: {
    fontSize: 20,
    fontWeight: 'bold',
    color: 'black',
  },
  menuItem: {
    paddingVertical: 15,
  },
  menuText: {
    fontSize: 16,
    color: 'black',
  },
  divider: {
    height: 1,
    backgroundColor: '#eee',
    marginVertical: 10,
  },
  sectionHeader: {
    fontSize: 14,
    color: 'gray',
    marginBottom: 10,
  },
  appearanceRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 10,
  },
  appearanceBtn: {
    padding: 10,
    borderWidth: 1,
    borderColor: '#eee',
    borderRadius: 8,
    flex: 1,
    alignItems: 'center',
    marginHorizontal: 4,
  },
  btnText: {
    color: 'black',
  },
  activeAppearance: {
    backgroundColor: '#e3f2fd',
    borderColor: '#2196F3',
  },
  activeText: {
    color: '#2196F3',
    fontWeight: 'bold',
  },
  logoutText: {
    color: 'red',
    fontWeight: 'bold',
  },
  closeButton: {
    marginTop: 20,
    backgroundColor: '#eee',
    padding: 15,
    borderRadius: 10,
    alignItems: 'center',
  },
  closeButtonText: {
    color: 'black',
    fontWeight: 'bold',
  },
});
