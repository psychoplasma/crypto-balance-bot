import { useEffect, useState } from 'react';
import Layout from '../../components/Layout';
import { deleteSubscription, getSubscriptions } from '../../lib/api';
import { Subscription } from '../../lib/types';
import { useAuth } from '../../context/AuthContext';
import { redirect } from 'next/navigation';
import SubscriptionCard from '../../components/SubscriptionCard';
import './user.css';

const UserPage = async (props: { userId: string }) => {
  // const [loading, setLoading] = useState(false);
  const loading = false;
  const { isAuthenticated } = useAuth();

  if (!isAuthenticated) {
    redirect('/login');
  }

  const fetchSubscriptions = async (userId: string) => {
    try {
      const subs = await getSubscriptions(userId);
      console.dir('user subs: ', subs);
    } catch (error) {
      console.error('Error fetching subscriptions:', error);
    } finally {
      // setLoading(false);
    }
  };

  let subscriptions: Subscription[];
  if (props.userId) {
    subscriptions = await fetchSubscriptions(props.userId);
  }

  const handleUnsubscribe = async (id: string, currency: string, address: string) => {
    try {
      await deleteSubscription(id, currency, address);
      // setSubscriptions((subscriptions as Subscription[]).filter(
      //   sub => sub.currency !== currency && sub.account !== address) as never,
      // );
    } catch (error) {
      console.error('Error unsubscribing:', error);
    }
  };

  return (
    <Layout>
      <div className="user-container">
        <h1>Your Subscriptions</h1>
        {loading ? (
          <div className="loading">Loading your subscriptions...</div>
        ) : (
          <div className="subscriptions-list">
            {subscriptions.map(subscription => (
              <SubscriptionCard
                key={subscription.id}
                subscription={subscription}
                onUnsubscribe={handleUnsubscribe}
              />
            ))}
          </div>
        )}
      </div>
    </Layout>
  );
};

export default UserPage;