package client

// Controller stops client in case stop signal received via ClientDoneChan.
func (c *Client) controller() {
	for d := range c.ClientDoneChan {
		if d {
			c.close(false)

			return
		}
	}
}
