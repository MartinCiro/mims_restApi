import { connect, NatsConnection, JSONCodec, Subscription } from 'nats';

class NatsService {
  private nc: NatsConnection | null = null;
  private readonly servers: string = 'nats://nats_server:4222';
  private readonly codec = JSONCodec();

  async connect(): Promise<NatsConnection> {
    if (!this.nc) {
      try {
        this.nc = await connect({ servers: this.servers });

        // Manejo de eventos de conexión
        this.nc.closed().then((err) => {
          this.nc = null;
        });
      } catch (error) {
        throw error;
      }
    }
    return this.nc;
  }

  async publish<T>(subject: string, message: T): Promise<void> {
    const nc = await this.connect();
    nc.publish(subject, this.codec.encode(message));
  }

  async request<T, R>(subject: string, message: T, timeout = 2000): Promise<R> {
    const nc = await this.connect();
    const response = await nc.request(subject, this.codec.encode(message), { timeout });
    return this.codec.decode(response.data) as R;
  }

  async subscribe<T>(subject: string, callback: (data: T) => void): Promise<Subscription> {
    const nc = await this.connect();
    const sub = nc.subscribe(subject);
    
    (async () => {
      for await (const msg of sub) {
        callback(this.codec.decode(msg.data) as T);
      }
    })();

    return sub;
  }
}

export const natsService = new NatsService();
