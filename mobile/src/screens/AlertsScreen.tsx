import React, { useState, useLayoutEffect } from 'react';
import { View, Text, FlatList, StyleSheet, ActivityIndicator, TouchableOpacity, TextInput, Modal, Switch } from 'react-native';
import { useQuery, gql } from 'urql';
import { useNavigation } from '@react-navigation/native';
import UserProfileMenu from '../components/UserProfileMenu';

const ALERTS_QUERY = gql`
  query alertsList($input: AlertSearchOptions) {
    alerts(input: $input) {
      nodes {
        id
        alertID
        status
        summary
        service {
          name
        }
      }
    }
  }
`;

type FilterType = 'active' | 'unacknowledged' | 'acknowledged' | 'closed' | 'all';

import { useTheme } from '../context/ThemeContext';

export default function AlertsScreen() {
  const navigation = useNavigation();
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState<FilterType>('active');
  const [favoritesOnly, setFavoritesOnly] = useState(true);
  const [showFilter, setShowFilter] = useState(false);
  const [showSearch, setShowSearch] = useState(false);
  const [showProfile, setShowProfile] = useState(false);

  const { colors } = useTheme();

  const getStatusFilter = (s: FilterType) => {
    switch (s) {
      case 'acknowledged': return ['StatusAcknowledged'];
      case 'unacknowledged': return ['StatusUnacknowledged'];
      case 'closed': return ['StatusClosed'];
      case 'all': return [];
      default: return ['StatusAcknowledged', 'StatusUnacknowledged'];
    }
  };

  const [result] = useQuery({
    query: ALERTS_QUERY,
    variables: {
      input: {
        filterByStatus: getStatusFilter(filter),
        first: 50,
        search: search,
        favoritesOnly: favoritesOnly,
      },
    },
  });


  useLayoutEffect(() => {
    navigation.setOptions({
      headerRight: () => (
        <View style={{ flexDirection: 'row', alignItems: 'center' }}>
          <TouchableOpacity onPress={() => setShowSearch(!showSearch)} style={{ marginRight: 15 }}>
            <Text style={{ fontSize: 18 }}>🔍</Text>
          </TouchableOpacity>
          <TouchableOpacity onPress={() => { setShowFilter(true); setShowProfile(false); }} style={{ marginRight: 15 }}>
            <Text style={{ fontSize: 18 }}>⚙️</Text>
          </TouchableOpacity>
          <TouchableOpacity onPress={() => { setShowProfile(true); setShowFilter(false); }} style={{ marginRight: 15 }}>
            <Text style={{ fontSize: 24, color: 'white' }}>👤</Text>
          </TouchableOpacity>
        </View>
      ),
    });
  }, [navigation, showSearch]);

  const { data, fetching, error } = result;

  const renderHeader = () => (
    <View>
      {showSearch && (
        <TextInput
          style={[styles.searchInput, { backgroundColor: colors.border, color: colors.text }]}
          placeholder="Search alerts..."
          placeholderTextColor={colors.subText}
          value={search}
          onChangeText={setSearch}
        />
      )}
      <View style={[styles.infoBanner, { backgroundColor: colors.primary + '20' }]}>
        <Text style={[styles.infoText, { color: colors.primary }]}>
          {favoritesOnly
            ? `Showing ${filter} alerts you are on-call for and from favorites.`
            : `Showing ${filter} alerts for all services.`}
        </Text>
      </View>
    </View>
  );

  if (fetching && !data) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" />
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.center}>
        <Text style={styles.error}>Error: {error.message}</Text>
      </View>
    );
  }

  const alerts = data?.alerts?.nodes || [];

  return (
    <View style={[styles.container, { backgroundColor: colors.background }]}>
      {renderHeader()}
      <FlatList
        data={alerts}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <View style={[styles.item, { borderBottomColor: colors.border, borderLeftColor: item.status === 'StatusUnacknowledged' ? 'red' : item.status === 'StatusAcknowledged' ? 'orange' : 'green', borderLeftWidth: 4 }]}>
            <Text style={[styles.title, { color: colors.text }]}>#{item.alertID}: {item.status.replace('Status', '')}</Text>
            <Text style={[styles.summary, { color: colors.text }]}>{item.summary}</Text>
            <Text style={[styles.service, { color: colors.subText }]}>{item.service.name}</Text>
          </View>
        )}
        ListEmptyComponent={<Text style={[styles.empty, { color: colors.subText }]}>No results</Text>}
      />

      {showFilter && (
        <View style={styles.modalOverlay}>
          <TouchableOpacity style={styles.overlayBackground} activeOpacity={1} onPress={() => setShowFilter(false)} />
          <View style={[styles.modalContent, { backgroundColor: colors.card }]}>
            <Text style={[styles.modalTitle, { color: colors.text }]}>Filter Alerts</Text>

            <View style={styles.filterRow}>
              <Text style={{ color: colors.text }}>Favorites Only</Text>
              <Switch value={favoritesOnly} onValueChange={setFavoritesOnly} />
            </View>

            <Text style={[styles.sectionTitle, { color: colors.text }]}>Status</Text>
            {['active', 'unacknowledged', 'acknowledged', 'closed', 'all'].map((s) => (
              <TouchableOpacity key={s} onPress={() => setFilter(s as FilterType)} style={[styles.filterOption, { borderBottomColor: colors.border }]}>
                <Text style={{ fontWeight: filter === s ? 'bold' : 'normal', color: filter === s ? colors.primary : colors.text }}>
                  {s.charAt(0).toUpperCase() + s.slice(1)}
                </Text>
              </TouchableOpacity>
            ))}

            <TouchableOpacity onPress={() => setShowFilter(false)} style={styles.closeButton}>
              <Text style={styles.closeButtonText}>Done</Text>
            </TouchableOpacity>
          </View>
        </View>
      )}

      <UserProfileMenu visible={showProfile} onClose={() => setShowProfile(false)} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  error: {
    color: 'red',
    padding: 20,
  },
  item: {
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  title: {
    fontWeight: 'bold',
    marginBottom: 4,
  },
  summary: {
    fontSize: 16,
    marginBottom: 4,
  },
  service: {
    color: 'gray',
    fontSize: 12,
  },
  searchInput: {
    padding: 10,
    backgroundColor: '#f0f0f0',
    margin: 10,
    borderRadius: 5,
  },
  infoBanner: {
    padding: 10,
    backgroundColor: '#e3f2fd',
    marginBottom: 5,
  },
  infoText: {
    color: '#0d47a1',
    fontStyle: 'italic',
    textAlign: 'center',
  },
  empty: {
    textAlign: 'center',
    marginTop: 20,
    color: 'gray',
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
  modalContent: {
    backgroundColor: 'white',
    padding: 20,
    borderTopLeftRadius: 20,
    borderTopRightRadius: 20,
  },
  modalTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 20,
    textAlign: 'center',
  },
  filterRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 20,
  },
  sectionTitle: {
    fontWeight: 'bold',
    marginTop: 10,
    marginBottom: 10,
  },
  filterOption: {
    paddingVertical: 10,
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  closeButton: {
    marginTop: 20,
    backgroundColor: '#2196F3',
    padding: 15,
    borderRadius: 5,
    alignItems: 'center',
  },
  closeButtonText: {
    color: 'white',
    fontWeight: 'bold',
  },
});
