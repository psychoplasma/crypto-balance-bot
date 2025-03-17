import { Subscription } from '../lib/types';
import './SubscriptionCard.css';

interface SubscriptionCardProps {
    subscription: Subscription;
    onUnsubscribe: (id: string, currency: string, address: string) => Promise<void>;
}

const SubscriptionCard = ({ subscription, onUnsubscribe }: SubscriptionCardProps) => {
  return (
    <div className="subscription-card">
      <div className="subscription-info">
        <span className="currency">{subscription.currency}</span>
        <span className="address">{subscription.account}</span>
      </div>
      <button 
        className="unsubscribe-btn"
        onClick={() => onUnsubscribe(subscription.id, subscription.currency, subscription.account)}
      >
        Unsubscribe
      </button>
    </div>
  );
};

export default SubscriptionCard;