import { useState } from 'react';

interface Currency {
  symbol: string;
  name: string;
}

interface SubscriptionFormProps {
  currencies: Currency[];
  onSubmit: (formData: FormData) => Promise<void>;
}

const SubscriptionForm = ({ currencies, onSubmit }: SubscriptionFormProps) => {
  const [newSubscription, setNewSubscription] = useState({
    currency: '',
    address: ''
  });

  return (
    <form className="subscription-form" action={onSubmit}>
      <select
        name="currency"
        value={newSubscription.currency}
        onChange={(e) => setNewSubscription({
          ...newSubscription,
          currency: e.target.value
        })}
        required
      >
        <option value="">Select Currency</option>
        {currencies.map(currency =>
          <option key={currency.symbol} value={currency.symbol}>{currency.name}</option>
        )}
      </select>
      <input
        type="text"
        name="address"
        placeholder="Enter wallet address"
        value={newSubscription.address}
        onChange={(e) => setNewSubscription({
          ...newSubscription,
          address: e.target.value
        })}
        required
      />
      <button type="submit">Create Subscription</button>
    </form>
  );
};

export default SubscriptionForm;
