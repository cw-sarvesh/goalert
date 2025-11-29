import React from 'react';
import { View, Text, StyleSheet } from 'react-native';

interface Props {
  name: string;
}

const PlaceholderScreen: React.FC<Props> = ({ name }) => {
  return (
    <View style={styles.container}>
      <Text style={styles.text}>{name} Screen</Text>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  text: {
    fontSize: 20,
    fontWeight: 'bold',
  },
});

export default PlaceholderScreen;
