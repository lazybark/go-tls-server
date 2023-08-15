package server

func (s *Server) Stop() error {
	s.cancel()
	s.listener.Close()

	var err error
	s.connPoolMutex.Lock()
	for _, conn := range s.connPool {
		err = conn.Close()
		if err != nil && !s.sConfig.SuppressErrors {
			s.ErrChan <- err
		}
	}
	s.connPoolMutex.Unlock()

	//At this point no routine will be left that can write in these channels
	close(s.ConnChan)
	close(s.ErrChan)

	return nil
}
