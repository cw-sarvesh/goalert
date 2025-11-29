import React from 'react';
import { createDrawerNavigator } from '@react-navigation/drawer';
import {
  AlertsScreen,
  RotationsScreen,
  SchedulesScreen,
  EscalationPoliciesScreen,
  ServicesScreen,
  UsersScreen,
  AdminScreen,
  WizardScreen,
} from '../screens';

const Drawer = createDrawerNavigator();

const DrawerNavigator = () => {
  return (
    <Drawer.Navigator initialRouteName="Alerts">
      <Drawer.Screen name="Alerts" component={AlertsScreen} />
      <Drawer.Screen name="Rotations" component={RotationsScreen} />
      <Drawer.Screen name="Schedules" component={SchedulesScreen} />
      <Drawer.Screen name="Escalation Policies" component={EscalationPoliciesScreen} />
      <Drawer.Screen name="Services" component={ServicesScreen} />
      <Drawer.Screen name="Users" component={UsersScreen} />
      <Drawer.Screen name="Admin" component={AdminScreen} />
      <Drawer.Screen name="Wizard" component={WizardScreen} />
    </Drawer.Navigator>
  );
};

export default DrawerNavigator;
