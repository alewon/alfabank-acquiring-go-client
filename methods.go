package alfabank

import "context"

// Register calls register.do.
func (c *Client) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	var resp RegisterResponse
	if err := c.postForm(ctx, "register.do", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RegisterPreAuth calls registerPreAuth.do.
func (c *Client) RegisterPreAuth(ctx context.Context, req RegisterPreAuthRequest) (*RegisterPreAuthResponse, error) {
	var resp RegisterPreAuthResponse
	if err := c.postForm(ctx, "registerPreAuth.do", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOrderStatusExtended calls getOrderStatusExtended.do.
func (c *Client) GetOrderStatusExtended(ctx context.Context, req GetOrderStatusExtendedRequest) (*GetOrderStatusExtendedResponse, error) {
	var resp GetOrderStatusExtendedResponse
	if err := c.postForm(ctx, "getOrderStatusExtended.do", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Deposit calls deposit.do.
func (c *Client) Deposit(ctx context.Context, req DepositRequest) (*OperationResponse, error) {
	var resp OperationResponse
	if err := c.postForm(ctx, "deposit.do", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Reverse calls reverse.do.
func (c *Client) Reverse(ctx context.Context, req ReverseRequest) (*OperationResponse, error) {
	var resp OperationResponse
	if err := c.postForm(ctx, "reverse.do", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Refund calls refund.do.
func (c *Client) Refund(ctx context.Context, req RefundRequest) (*OperationResponse, error) {
	var resp OperationResponse
	if err := c.postForm(ctx, "refund.do", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Decline calls decline.do.
func (c *Client) Decline(ctx context.Context, req DeclineRequest) (*OperationResponse, error) {
	var resp OperationResponse
	if err := c.postForm(ctx, "decline.do", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ProcessRawSumRefund calls processRawSumRefund.do.
func (c *Client) ProcessRawSumRefund(ctx context.Context, req ProcessRawSumRefundRequest) (*OperationResponse, error) {
	var resp OperationResponse
	if err := c.postForm(ctx, "processRawSumRefund.do", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ProcessRawPositionRefund calls processRawPositionRefund.do.
func (c *Client) ProcessRawPositionRefund(ctx context.Context, req ProcessRawPositionRefundRequest) (*OperationResponse, error) {
	var resp OperationResponse
	if err := c.postForm(ctx, "processRawPositionRefund.do", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}