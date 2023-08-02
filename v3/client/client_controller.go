package client

// Controller stops client in case stop signal recieved
func (c *Client) Controller() {
	for {
		select {
		case d := <-c.ClientDoneChan:
			if d {
				c.close(false)
				return
			}
		case <-c.ctx.Done():
			c.close(false)
			return
		}

	}
}
