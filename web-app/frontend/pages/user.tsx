import { useEffect, useState } from 'react';
import Layout from '../components/Layout';
import { deleteSubscription, getSubscriptions } from '../api/api';
import { Subscription, User } from '../api/types';

const UserPage = (props: { userId: string }) => {
    // const [subscriptions, setSubscriptions] = useState<Subscription[]>([
    //     { id: '1', currency: 'BTC', account: '1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa',
    //         userId: 'user1',
    //         blockHeight: 100,
    //         startingBlockHeight: 0,
    //         filters: '',
    //     },
    //     { id: '2', currency: 'ETH', account: '0x742d35Cc6634C0532925a3b844Bc454e4438f44e',
    //         userId: 'user1',
    //         blockHeight: 100,
    //         startingBlockHeight: 0,
    //         filters: '',
    //     },
    //     { id: '3', currency: 'DOGE', account: 'DH5yaieqoZN36fDVciNyRueRGvGLR3mr7L',
    //         userId: 'user1',
    //         blockHeight: 100,
    //         startingBlockHeight: 0,
    //         filters: '',
    //     }
    // ]);
    const [subscriptions, setSubscriptions] = useState<Subscription[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        const fetchSubscriptions = async (userId: string) => {
            try {
                const subs = await getSubscriptions(userId);
                setSubscriptions(subs as never);
            } catch (error) {
                console.error('Error fetching subscriptions:', error);
            } finally {
                setLoading(false);
            }
        };

        if (props.userId) {
            fetchSubscriptions(props.userId);
        }
    }, []);

    const handleUnsubscribe = async (id: string, currency: string, address: string) => {
        try {
            await deleteSubscription(id, currency, address);
            setSubscriptions((subscriptions as Subscription[]).filter(
                sub => sub.currency !== currency && sub.account !== address) as never,
            );
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
                        {subscriptions.map(s => (
                            <div key={s.id} className="subscription-card">
                                <div className="subscription-info">
                                    <span className="currency">{s.currency}</span>
                                    <span className="address">{s.account}</span>
                                </div>
                                <button 
                                    className="unsubscribe-btn"
                                    onClick={() => handleUnsubscribe(s.id, s.currency, s.account)}
                                >
                                    Unsubscribe
                                </button>
                            </div>
                        ))}
                    </div>
                )}
                <style jsx>{`
                    .user-container {
                        max-width: 900px;
                        margin: 2rem auto;
                        padding: 0 1rem;
                    }

                    h1 {
                        color: #333;
                        margin-bottom: 2rem;
                        text-align: center;
                    }

                    .loading {
                        text-align: center;
                        color: #666;
                        padding: 2rem;
                    }

                    .subscriptions-list {
                        display: grid;
                        gap: 1rem;
                    }

                    .subscription-card {
                        display: flex;
                        justify-content: space-between;
                        align-items: center;
                        background: white;
                        padding: 1.5rem;
                        border-radius: 8px;
                        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
                        transition: transform 0.2s ease;
                    }

                    .subscription-card:hover {
                        transform: translateY(-2px);
                    }

                    .subscription-info {
                        display: grid;
                        gap: 0.5rem;
                    }

                    .currency {
                        font-weight: bold;
                        color: #0070f3;
                    }

                    .address {
                        color: #666;
                        font-family: monospace;
                    }

                    .balance {
                        color: #2ecc71;
                        font-weight: 500;
                    }

                    .unsubscribe-btn {
                        padding: 0.5rem 1rem;
                        background-color: #ff4757;
                        color: white;
                        border: none;
                        border-radius: 4px;
                        cursor: pointer;
                        transition: background-color 0.2s ease;
                    }

                    .unsubscribe-btn:hover {
                        background-color: #ff6b81;
                    }

                    @media (max-width: 600px) {
                        .subscription-card {
                            flex-direction: column;
                            gap: 1rem;
                            text-align: center;
                        }
                    }
                `}</style>
            </div>
        </Layout>
    );
};

export default UserPage;