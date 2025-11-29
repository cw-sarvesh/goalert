import React, { createContext, useState, useContext, useEffect } from 'react';
import { useColorScheme } from 'react-native';

type ThemeType = 'light' | 'dark' | 'system';

interface ThemeColors {
  background: string;
  text: string;
  subText: string;
  border: string;
  card: string;
  primary: string;
  error: string;
  success: string;
  warning: string;
  modalOverlay: string;
}

const lightColors: ThemeColors = {
  background: '#ffffff',
  text: '#000000',
  subText: 'gray',
  border: '#eeeeee',
  card: '#ffffff',
  primary: '#2196F3',
  error: 'red',
  success: 'green',
  warning: 'orange',
  modalOverlay: 'rgba(0,0,0,0.5)',
};

const darkColors: ThemeColors = {
  background: '#1F1F1F', // Greyish background
  text: '#ffffff',
  subText: '#aaaaaa',
  border: '#333333',
  card: '#2C2C2C', // Slightly lighter grey for cards
  primary: '#90caf9',
  error: '#ef5350',
  success: '#66bb6a',
  warning: '#ffa726',
  modalOverlay: 'rgba(255,255,255,0.1)',
};

interface ThemeContextType {
  theme: ThemeType;
  setTheme: (theme: ThemeType) => void;
  colors: ThemeColors;
  isDark: boolean;
}

const ThemeContext = createContext<ThemeContextType>({
  theme: 'system',
  setTheme: () => {},
  colors: lightColors,
  isDark: false,
});

export const ThemeProvider = ({ children }: { children: React.ReactNode }) => {
  const systemScheme = useColorScheme();
  const [theme, setTheme] = useState<ThemeType>('system');

  const isDark = theme === 'system' ? systemScheme === 'dark' : theme === 'dark';
  const colors = isDark ? darkColors : lightColors;

  return (
    <ThemeContext.Provider value={{ theme, setTheme, colors, isDark }}>
      {children}
    </ThemeContext.Provider>
  );
};

export const useTheme = () => useContext(ThemeContext);
