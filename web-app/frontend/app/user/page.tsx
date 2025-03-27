'use client';

import { useEffect, useState } from 'react';
import Layout from '@/components/Layout';
import { Subscription } from '@/lib/types';
import SubscriptionCard from '@/components/SubscriptionCard';
import SubscriptionForm from '@/components/SubscriptionForm';
import {
  createSubscriptionAction,
  deleteSubscriptionAction,
  getSubscriptionsAction,
} from '@/actions/actions';
import './user.css';
import ProtectedRoute from '@/components/ProtectedRoute';

const CURRENCIES = [
  {
    symbol: 'BTC',
    name: 'Bitcoin',
  },
  {
    symbol: 'ETH',
    name: 'Ethereum',
  }
];

const UserPage = () => {
  const [loading, setLoading] = useState(false);
  const [subscriptions, setSubscriptions] = useState<Subscription[]>([]);

  const fetchSubscriptions = async () => {
    try {
      const subs = await getSubscriptionsAction();
      setSubscriptions(subs);
    } catch (error) {
      console.error('Error fetching subscriptions:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchSubscriptions() }, []);

  const handleUnsubscribe = async (currency: string, address: string) => {
    try {
      await deleteSubscriptionAction(currency, address);
      setSubscriptions(subscriptions.filter(
        sub => sub.currency !== currency && sub.account !== address),
      );
    } catch (error) {
      console.error('Error unsubscribing:', error);
    }
  };

  const handleCreateSubscription = async (formData: FormData) => {
    try {
      const newSub = await createSubscriptionAction(
        formData.get('currency') as string,
        formData.get('address') as string,
        0,
        0,
      );
      setSubscriptions([...subscriptions, newSub]);
    } catch (error) {
      console.error('Error creating subscription:', error);
    }
  };

  return (
    <ProtectedRoute>
      <Layout>
        <div className="user-container">
          <div className="header-section">
            <h1>Create Subscription</h1>
          </div>

          <SubscriptionForm onSubmit={handleCreateSubscription} currencies={CURRENCIES} />

          <div className="header-section">
            <h1>Your Subscriptions</h1>
          </div>

          {loading ? (
            <div className="loading">Loading your subscriptions...</div>
          ) : (
            <div className="subscriptions-list">
              {subscriptions.map(sub => (
                <SubscriptionCard
                  key={sub.id}
                  subscription={sub}
                  onUnsubscribe={handleUnsubscribe}
                />
              ))}
            </div>
          )}
        </div>
      </Layout>
    </ProtectedRoute>
  );
};

export default UserPage;
