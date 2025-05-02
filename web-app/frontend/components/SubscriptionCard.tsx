import { formatCurrency } from '@/lib/utils';
import { Subscription } from '../lib/types';
import './SubscriptionCard.css';

interface SubscriptionCardProps {
    subscription: Subscription;
    onUnsubscribe: (currency: string, address: string) => Promise<void>;
}

const SubscriptionCard = ({ subscription, onUnsubscribe }: SubscriptionCardProps) => {
  const spent = formatCurrency(
    BigInt(subscription.totalSpent),
    BigInt(subscription.decimal),
    subscription.currency.toUpperCase(),
  );
  const received = formatCurrency(
    BigInt(subscription.totalReceived),
    BigInt(subscription.decimal),
    subscription.currency.toUpperCase(),
  );
  const netBalance = formatCurrency(
    BigInt(subscription.totalReceived) - BigInt(subscription.totalSpent),
    BigInt(subscription.decimal),
    subscription.currency.toUpperCase(),
  );
  
  return (
    <div className="subscription-card">
      <div className="card-header">
        <div className="currency-symbol">{subscription.currency.toUpperCase()}</div>
        <div className="wallet-address">{subscription.account}</div>
      </div>

      <div className="card-content">
        <div className="stat-group">
          <div className="stat-item">
            <div className="stat-label">Received</div>
            <div className="stat-value">
              {received} {subscription.currency}
            </div>
          </div>

          <div className="stat-item">
            <div className="stat-label">Spent</div>
            <div className="stat-value">
              {spent} {subscription.currency}
            </div>
          </div>

          <div className="stat-item">
            <div className="stat-label">Net Balance</div>
            <div className="stat-value">
              {netBalance} {subscription.currency}
            </div>
          </div>

          <div className="stat-item">
            <div className="stat-label">Last update</div>
            <div className="stat-value">Block#{subscription.blockHeight}</div>
          </div>
        </div>

        <button 
          className="unsubscribe-btn"
          onClick={() => onUnsubscribe(subscription.currency, subscription.account)}
        >
          Unsubscribe
        </button>
      </div>
    </div>
  );
};

export default SubscriptionCard;
